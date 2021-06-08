# Channels

The purpose of this example is to demonstrate communication using channels,
and the select statement.

## Design
              __--> worker >--__
             /                  \
            /-----> worker >--___\
    source /                      \ destination
            \-----> worker >-----/
             \                  /
              \_--> worker >--_/

The premise here is to get a feel of how Go channels behave, and how the
`select` statement works with channels. The data flows from a source, and
is then multiplexed among the channels.  Each channel has a worker
routine to accept the data, and then pass it on to its output. Of course
in practice the worker would do something with the data, but here I
just want to work with the flow.

The demultiplexing worker selects among the channels, and aggregates the
data back to a single channel.  When all of the channels are closed by
their workers, this routine will exit.

The final destination is a channel where the main routine is waiting for
data.  Once the demuxChan closes, the main routine will exit.

## Things Learned

Channels behave differently on a single-core vs multi-core machine.  On an
AWS EC2 instance, using a buffer value higher than 16 caused a slow-down.
On a multi-core machine, no slow-down is noted.

