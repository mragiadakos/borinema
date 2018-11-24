#!/bin/bash

# To download dependencies, we have to download them with the GO111MODULE=off flag
# The reason is that the gopherjs does not read the vendor folder yet for projects that use go.mod
# So we have to download the dependencies in the GOPATH
GO111MODULE=off gopherjs get github.com/HuckRidgeSW/hvue
GO111MODULE=off gopherjs get github.com/gopherjs/gopherwasm/js
GO111MODULE=off gopherjs get honnef.co/go/js/xhr
GO111MODULE=off gopherjs get github.com/cathalgarvey/fmtless
gopherjs build -m .