# IP2ASN

A command-line tool to look up the ASN (Autonomous System Number) of one or more IPv4 or IPv6 addresses.

The tool uses the WHOIS service provided by Cymru to look up the ASN for each IP address.

## Installation

To install `ip2asn`, run the following command:

```
go install github.com/melvinsh/ip2asn@latest
```

## Usage

The tool reads IP addresses from standard input, one address per line. For example:

```
$ echo "8.8.8.8" | ip2asn
15169
```

You can also pipe in a file containing a list of IP addresses:

```
cat ips.txt | ip2asn
```

The output is the ASN number for each IP address. If the tool is unable to look up the ASN for an IP address, it will print an error message to standard error and move on to the next address.
