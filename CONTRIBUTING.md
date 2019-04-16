If you'd like to contribute, please fork this repository and publish PRs for review. If you'd like to get in touch with the author, email ```dabdine <at> rapid7.com```. We are always looking for further conversion of capabilities supported in the ruby-based DAP, as well as general performance improvements to make godap as fast as possible.

# Testing

There are two testing frameworks in place.

Code-level tests are required, and must be written in [goconvey](https://github.com/smartystreets/goconvey).
This library allows behavior-driven development and testing. It is also compatible with standard golaing 
`testing.T`.


Additionally, [bats](https://github.com/sstephenson/bats) is currently used to run integration
tests.  [travis-ci](https://travis-ci.com) will automatically run all `bats` tests defined in this project 
upon each PR.  You are encouraged to add tests as you add/convert functionality from ruby-based dap to make
the port easier.

To run tests outside of travis-ci:

```
docker build -t godap_bats -f Dockerfile.bats_testing . &&  docker run --rm -it --name godap_bats -v "$PWD":/opt/bats_testing godap_bats
```
