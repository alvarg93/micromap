# Micromap (WIP)

Micromap is a suite of command line tools to generate a visual representation of your application dependencies.

![Example graph](micromap.png?raw=true "Micromap graph")

## Dependencies

This package depends heavily on [Graphviz](https://www.graphviz.org/).

## Commands

1. `mmap` - (micromap) generate a Dot file and an image file from a yaml specification.

2. `mgen` - (microgen) generate a yaml specification file from annotations in code comments

3. `mmerge` - (micromerge) additional wrapper over `m4` that supports URLs as well. Will add a more sophisticated option of aggregating yamls at a later time.

## Motivation

Honestly, I just wanted something that would make it easier to keep microservice network diagrams up to date.

## What's the point?

1. As I mentioned - it is really difficult to keep network diagrams up to date in a microservice architecture. Nobody wants to draw it from scratch every single time.

2. Coupling! Graphs are a great visual tool to detect heavy coupling and 
plan some refactor of your code or architecture.

# Usage

## Yaml

Write your specification in a `.yml` file.

```yaml
# ./micromap.yml
config:
  app: MyApp
groups:
- name: DataCenter
  deps:
  - name: Postgres
    typ: db
    parent: DataCenter
- name: AWS
deps:
- name: SNS
  typ: queue
  parent: AWS
  rels:
  - path: NotifyUser
    dir: forward
- name: SQS
  typ: queue
  parent: AWS
- name: PostgresBackup
  typ: db
  parent: DataCenter
rels:
- service: Postgres
  path: Orders
- service: Postgres
  path: Products
- service: SQS
  path: OutOfStock
  dir: back
```

Then generate the graph.

    $ mmap -y=micromap.yml -d=micromap.dot -i=micromap.png

## Generate from comments

You can also write it in comments if it helps you keep it up to date:

```javascript
    //@micromap
    //config:
    //  app: MyApp
    //groups:
    //  - name: DataCenter
    //    deps:
    //      - name: Postgres
    //        typ: db
    //        rels:
    //          - path: Orders
    //          - path: Products
    //  - name: AWS
    //    deps:
    //      - name: SNS
    //        typ: queue
    //        rels:
    //          - path: NotifyUser
    //            dir: forward
    //      - name: SQS
    //        typ: queue
    //deps:
    //  - name: Postgres2
    //    typ: db
    //    parent: DataCenter
    //rels:
    //  - path: OutOfStock
    //    dir: back
    //    service: SQS
```

Then generate a yaml file.

    $ mgen -y=micromap.yml -x=.go -r

## Multiple dot files (yamls coming soon...)

Large projects = multiple services = separate specifications. And that's okay!

Merge graphs from different projects using `mmerge`.

    $ mmerge -d=merged.dot -i=merged.png micromap.dot ../micromap2.dot http://example.com/micromap.dot