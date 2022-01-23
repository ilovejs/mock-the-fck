# Mock the fck

Originally it's `Mockery-example`

Example case for [mockery issue #128](https://github.com/vektra/mockery/issues/128) filed with the golang tool "[mockery](https://github.com/vektra/mockery)".

But right now, I added more complex examples

#### Topics

- Mock
- 
#### Scripts

- Remove duplicated vscode plugin on mac

#### Mock notes:

##### How to generate S3API from vendor code ?

use `go mod vendor` to download code

*warning*

- s3 is a struct, so

`
mockgen --build_flags=--mod=mod github.com/aws/aws-sdk-go/service/s3 S3
`
is wrong

##### 1. Mockgen 

[S3API is interface](https://github.com/aws/aws-sdk-go/blob/main/service/s3/s3iface/interface.go)


`
mockgen --build_flags=--mod=mod github.com/aws/aws-sdk-go/service/s3/s3iface S3API  > mocks/second_mocks.go
`

##### 2. Mockery

What if I want mockery api for mocking and generation ??

1. add `_ "github.com/aws/aws-sdk-go/service/s3/s3iface/interface.go"`
2. go mod vendor
3. run below, but mocks/S3API.go will be overwritten !!!

```shell
mockery --dir vendor/github.com/aws/aws-sdk-go/service/s3/s3iface --name S3API
```

##### other takeaway cmds

run `go mod vendor` to pull code

`mockgen --build_flags=--mod=mod vendor/github.com/aws/aws-sdk-go/service/s3 S3`

`import _ "github.com/golang/mock/mockgen/model"`

##### Issue of method 2 under Q&A readme of mockgen

Better try method 3 !

`mockgen vendor/github.com/aws/aws-sdk-go/service/s3 S3`

```
prog.go:14:2: use of vendored package not allowed
prog.go:12:2: no required module provides package github.com/golang/mock/mockgen/model: go.mod file not found in current directory or any parent directory; see 'go help modules'
prog.go:14:2: vendor/github.com/aws/aws-sdk-go/service/s3 must be imported as github.com/aws/aws-sdk-go/service/s3
2022/01/21 17:54:17 Loading input failed: exit status 1
```

##### Buggy Test Coverage

run test with coverage.
shows main.go 0% covered, but mocks folder has 33%!!!

Why on earth we calculate coverage stats on that bits !!

### Futures

1. Why gomock or mockery not generate from struct as extra feature,
as in case of s3 struct.

2. [gcp sdk](https://github.com/google/go-cloud/tree/master/internal/testing)
Has more good pratices

3. [more gcp](https://github.com/google/go-cloud/blob/master/internal/testing/setup/setup.go)
shows fake data has been set in google style.

But, this [blog](https://www.philosophicalhacker.com/2016/01/13/should-we-use-mocking-libraries-for-go-testing/) points out:

>   Google engineer prefer hand-code fake rather reflection mock

```txt
Andrew Gerrand mentions gomock, a mocking library, and reluctantly says

[mocking libraries like gomock] are fine, 
but I find that on balance the hand-written fakes tend be easier to reason about 
and clearer to see what's going on, but I'm not an enterprise go programmer 
so maybe people do need that so I don't know, but that's my advice.
```
4. [uber/kraken](https://github.com/uber/kraken/tree/master/mocks/lib/backend) 
Manually create s3 interface. But their Makefile mockgen heaps of stuff !!

5. [aws-go](https://github.com/aws/aws-sdk-go/blob/main/service/s3/s3iface/interface.go)
Lots of http code

6. [counterfeiter demo by mario](https://github.com/MarioCarrion/videos/tree/main/2020/11/24-go-tools-counterfeiter)
Probably interesting
