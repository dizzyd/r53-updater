r53-updater
===========

Simple utility to update a single A record in Route 53.

r53-updater uses [icanhazip.com](http://major.io/icanhazip-com-faq/) to get the
current public IP of your system and update a hosted zone on [Amazon's
Route53](http://aws.amazon.com/route53) DNS service.

Building
--------

    $ make

A successful build will leave an executable file, `r53-updater`. This executable
has no external dependencies, so you can safely copy this to another machine and
use it there.

Running
-------

You need to create a configuration file for the updater to use. This file needs
to contain the following information:

*  accesskey: AWS access key for the API call to Route53
*  secretkey: AWS secret key for the API call to Route53
*  zoneid: ID of the hosted zone to update on Route53
*  name: fully-qualified DNS name you wish to associate with the public IP
*  ttl: number of seconds in the A record's TTL

There is an example configuration file included with the source.

Once you have a config, run the updater:

    $ r53-updater -c example.config

If all goes well, nothing will be output. Any failure while print an error
message to the console.



