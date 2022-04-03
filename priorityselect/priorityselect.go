// Package priorityselect implementes a select-like construct that that receives messages from channels in
// order of user-specified priortiy.  When there are messages available on multiple channels, they will be
// delivered in priority order.
//
// TODO: define channel closing semantics
//
// TODO: clean shutdown
//
// TODO: think of a better name for this: if we are just providing access to the outgoing channel then the
// possibilities are far wider than just a `select` call.
package priorityselect

// PrioritisedChannel represents a channel configured for use in a ChannelGroup with a priority and a tag.
type PrioritisedChannel[A any] struct {
	Channel <-chan A
	// Priority specifies the priority of the channel within the group.  Higher value -> higher priority.
	Priority int
	// Tag is an arbitrary value used to identify the original source of a message from a ChannelGroup
	Tag any
}

// TaggedMessage
type TaggedMessage[A any] struct {
	Value A
	Tag   any
}

// ChannelGroup is a collection of channels with priorities and tags.
//
// Once a channel has been used as part of a ChannelGroup, it is inadvisable to receieve messages on it
// through any other means.  ChannelGroup maintains its own internal buffers so the results will be
// unpredictable.
type ChannelGroup[A any] struct {
	outgoing chan *TaggedMessage[A]
}

func NewChannelGroup[A any](channels ...*PrioritisedChannel[A]) *ChannelGroup[A] {
	// not pre-allocating any buffer space because we have no idea how much will be needed: possible
	// future enhancement would be to allow user to specify per-channel
	// TODO: need to be able to look up buffers by tag *and* iterate in priority order
	buffer := make(map[any][]A)

	// these buffer sizes were pulled out of the air, could do with a bit of thought
	intermediate := make(chan *TaggedMessage[A], len(channels)*4)
	outgoing := make(chan *TaggedMessage[A], len(channels)*4)

	for _, ch := range channels {
		// Collect values from incoming, tag and forward to intermediate.  One goroutine per channel is the best
		// way I've found to do this without spinning (so far).  This functionality is potentially useful on its
		// own so could be factored out into a separate package.
		//
		// TODO: there is a Select for an arbitrary number of channels in reflect(!)
		go func(ch *PrioritisedChannel[A]) {
			value, ok := <-ch.Channel
			if !ok { // channel was closed
				return
			}
			intermediate <- &TaggedMessage[A]{value, ch.Tag}
		}(ch)
	}

	// Collect values from intermediate and buffer until there are no more values waiting (receive would
	// block), then send out in priority order and repeat.
	go func() {
		var buffered int
		var msg *TaggedMessage[A]

		for {
			if buffered > 0 {
				select {
				case msg = <-intermediate:
				default:

				}
			} else { // nothing waiting to go out so we're ok to block
				msg = <-intermediate
			}

			if msg != nil {
				buffer[msg.Tag] = append(buffer[msg.Tag], msg.Value)
				msg = nil
			}
		}
	}()

	cg := &ChannelGroup[A]{intermediate, outgoing}

	return cg
}

func (channelGroup *ChannelGroup[A]) Outgoing() <-chan *TaggedMessage[A] {
	return channelGroup.outgoing
}
