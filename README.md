# U Up?

A small service to test whether your port is internet accessible.

Written in Go, in a fit of frustration.

## Limitations

### TCP 

The TCP checker returns success on an opened connection. It 
does not read or write any data to that connection.

### Quic

For the sake of security, various cryptographic properties are
not user-specified. Thus, the Quic checker does not verify TLS
certificates on connect. This also means the connect may fail
with `CRYPTO_ERROR` in various ways, as the host rejects
the arbitrary checker certificate that also lacks the correct 
application protocol names. Consequently, the Quic checker
returns success on an opened connection OR A `CRYPTO_ERROR`

### UDP Checker

The UDP checker sends up to k datagrams with message `U UP?` 
to the user-specified port. It assumes success if it receives any 
packet back. If the protocol does not reply to this message,
it still appears down, though it need not be.

## (Imagined) FAQs

### Why doesn't the JSON response return my IP address?

Because I don't want people using this service one that provides them with 
their IP address. Use someone else's [STUN](https://en.wikipedia.org/wiki/STUN)
service for that.

### Why can't I specify the IP address I'm curious about?

Because this isn't a general port scanner. It's a service to help you
write and debug your local network, especially if you are trying to 
write P2P applications (which I want you to do).

### Why can't I specify the payload?

Because this service exists only to verify the port is accessible. If 
you could specify the payload, it would provide a free service for 
anonymously executing attacks against the network you happen to be 
on, which might not be yours (evil Mallory!).

### Is this service logged?

Hrm.

