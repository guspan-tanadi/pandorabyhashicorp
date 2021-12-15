#!/bin/bash

function buildAndInstallDependencies {
    echo "Installing the Go SDK Generator into the GOBIN.."
    cd ./tools/generator-go-sdk
    go install .
    cd ../../

    echo "Building Wrapper.."
    cd ./tools/wrapper-go-sdk-generator
    go build -o wrapper-go-sdk-generator
    cd ../../
}

function runWrapper {
  local dataApiAssemblyPath=$1
  local outputDirectory=$2

  echo "Running Wrapper.."
  cd ./tools/wrapper-go-sdk-generator
  ./wrapper-go-sdk-generator \
    -data-api-assembly-path="../../$dataApiAssemblyPath"\
    -output-dir="../../$outputDirectory"
}

function prepareGoSdk {
  echo "TODO: conditionally checkout and prepare the Go SDK.."
}

function conditionallyCommitAndPushGoSdk {
  echo "TODO: conditionally commit/push the Go SDK"
}

function main {
  local dataApiAssemblyPath="data/Pandora.Api/bin/Debug/net6.0/Pandora.Api.dll"
  local outputDirectory="generated/" # TODO: replace with the submodule in time

  buildAndInstallDependencies
  prepareGoSdk
  runWrapper "$dataApiAssemblyPath" "$outputDirectory"
  conditionallyCommitAndPushGoSdk
}

main