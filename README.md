# arangoctl

Tool to declaratively manage ArangoDB collections (with indexes) and searchviews.

# Motivation

We are heavy user of Kafka. we manage our Kafka topics using [TopicCtl](https://github.com/segmentio/topicctl).

Similar to topicctl, we wanted to manage our Arango objects in a declarative way.

# Getting started

## Installation

You can set it up in one of the following ways:

1. go get github.com/psykidellic/arangoctl/cmd/arangoctl (tested with only Go 1.16)
2. Clone the repo and run: go build cmd/arangoctl

## Quick Tour

Setup a quick cluster using docker-compose (follow this blog) [spinning up arangodb using docker-compose](https://dev.to/sonyarianto/how-to-spin-arangodb-server-with-docker-and-docker-compose-3c00).

Apply the arango resources:

```
arangoctl apply --cluster-config examples/cluster-auth.yaml examples/resources/*
```

Now visit the Arango Web UI to see the two resources being present.

NOTE: Internally, we are using the [official Arango go-driver](https://github.com/arangodb/go-driver).

# Usage

## cluster-config

Every *arangoctl* command will have a common *--cluster-config* option that allows you to provide details of the Arango cluster on which the commands will be executed.

Currently, we only support connection with simple auth or no auth. Example of each is provided in the *examples* folder.

## Subcommnds

### apply [path to resource(s)]

Will declaratively manage resources. For collections, we create the collection if its not present and will also create any of the listed indexes. For existing collections, it will do a comparison of indexes between what is defined in the yaml and what is present and update accordingly (including adding/removing fields from indexes). Comparison is done by index and field name. It will DELETE indexes which are not defined in the resource specification.

The goal is that you will always use arangoctl and keep your resources defined in source control.

**NOTE:** As of now, we dont allow to delete resources. You will have to do it using other mechanism.

## Resource Definition

### Collections

See an example at *examples/resources/01-fooview.yaml*. You define a resource to create collection by setting "Collection" for Kind. meta.type can be one of:

1. edge
2. collection

As of now, the only spec for a collection is supported is spec.indexes. This defines a list of index, type and list of columns in the index.

### Searchviews

You can configure searchview index based [on the docs](https://www.arangodb.com/docs/stable/arangosearch-views.html). The same can be mapped to this driver structure [ArangoSearchViewProperties](https://github.com/arangodb/go-driver/blob/master/view_arangosearch.go#L235). The spec of arangoctl resource, follows the same structure.
