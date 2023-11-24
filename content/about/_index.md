---
title: "About the Dynamic Simulation Environment"
linkTitle: About
menu:
  main:
    weight: 10
---

{{% blocks/cover title="Dynamic Simulation Environment" image_anchor="bottom" height="auto" %}}

{{% /blocks/cover %}}


{{% blocks/lead %}}
The Dynamic Simulation Environment is a message based distributed simulation platform which defines the interface between models as signal vectors. Signal vectors are a logical grouping of either scalar or binary data which are exchanged between models at discrete points in time. Models may be developed in any programming language and may run on any operating system or hardware platform.

Simulations are constructed using models with a compositional approach, where system elements represented by the simulation are decomposed/recomposed from a collection of models according the broader simulation objectives. This compositional approach to simulation, called Functional Simulation (FSIL), supports the development of a variety of simulation strategies.

The Dynamic Simulation Environment (DSE) can be used to integrate legacy tools into a simulation, such integration techniques include; Gateway models which connect remote systems to a DSE simulation; and Model Compatibility Libraries which import foreign models into a DSE simulation. Several Gateway and Model Compatibility Libraries are available alongside the Model C Library, our C based implementation of the Dynamic Simulation Environment.
{{% /blocks/lead %}}


{{% blocks/section color="white" %}}

## Messages

The Dynamic Simulation Environment uses a distributed algorithm, based upon a small set of messages, which facilitate control and execution of simulations. Using those messages (and the associated distributed algorithm) a variety of simulation capabilities can be realised, including: Co-simulation, Federated (distributed) simulations, Cross-Platform simulations.

The Messages of the Dynamic Simulation Environment are used to:

* **Progress time in a Simulation:** the synchronisation algorithm which drives the Dynamic Simulation Environment is implemented with 3 messages; two of which are used to implement the Co-simulation algorithm, and a third message which is used for event driven progression of time.

* **Facilitate the exchange of Signals between Models:** signal values are first delta encoded (i.e. only changed values), and then embedded as a payload within messages. These messages are then exchanged with, and processed by, a Simulation Bus which implements the distributed algorithm of the Dynamic Simulation Environment.

* **Implementation of Model Interfaces:** the messages of the Dynamic Simulation Environment are defined using cross platform serialization libraries (Flatbuffers, MsgPack and gRPC). Model interfaces can therefore be written for any combination of operating System, architecture or programming language.

{{% /blocks/section %}}


{{% blocks/section color="white" %}}

## Signal Vectors

The Model C Library presents signals to a model developer as Signal Vectors; simple vector objects, either holding scalar or binary signals, and a collection of supporting methods. Typically, the only interactions a Model Developer will have with the Dynamic Simulation Environment will be related to Signal Vectors.

In the Dynamic Simulation Environment, when using the Model C Library, Signal Vectors are used to:

* **Provide a simple Model API:** only 2 functions are required to implement a Model in the Dynamic Simulation Environment.

* **Represent Signals:** all scalar types (bool, int8, float, double etc.) can represented in a scalar Signal Vector, and complex or streaming data types can be represented with a binary Signal Vector. A binary signal may also be annotated with a MIMEtype to further describe additional schemas or properties of the binary data.

* **Encapsulate Models from other modelling standards:** a Model Compatibility Layer may be implemented to extend the Model C Library to support a foreign modelling tool or framework. In that case, the native data types of the foreign model may be marshalled into either scalar or binary Signal Vectors as required, and then interact seamlessly with other models present in the simulation.

* **Connect to external simulations:** a Gateway Model allows remote/legacy simulation tools to be connected to a Dynamic Simulation Environment and exchange signals using Signal Vectors. Gateway models are written in the same programming language or framework as that of the remote/legacy simulation tool, and also run in that same environment.

{{% /blocks/section %}}


{{% blocks/section color="secondary" %}}

## The Core Platform

The Dynamic Simulation Environment Core Platform is a collection of schemas, libraries and tools which can be used to design and implement simulations using the compositional approach.


### Contributing

Want to contribute? Great! Each project repository has contribution instructions, however the general recommendation is to open an Issue or Pull Request. If your planning a larger contribution, get in contact first (e.g. open an Issue) and discuss the changes.

### More Information

Check out the docs, or open an issue on one of the repos.


{{% /blocks/section %}}
