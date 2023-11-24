---
title: Dynamic Simulation Environment
---

{{< blocks/cover title="" image_anchor="bottom" height="full" >}}
<a class="btn btn-lg btn-primary me-3 mb-4" href="about/">
  Learn More <i class="fas fa-arrow-alt-circle-right ms-2"></i>
</a>
<a class="btn btn-lg btn-secondary me-3 mb-4" href="https://github.com/boschglobal?q=dse&type=all">
  Git Hub <i class="fab fa-github ms-2 "></i>
</a>
{{< /blocks/cover >}}


{{% blocks/lead color="primary" %}}
The Dynamic Simulation Environment is a message based distributed simulation platform which defines the interface between models as signal vectors. Signal vectors are a logical grouping of either scalar or binary data which are exchanged between models at discrete points in time. Models may be developed in any programming language and may run on any operating system or hardware platform.

Simulations are constructed using models with a compositional approach, where system elements represented by the simulation are decomposed/recomposed from a collection of models according the broader simulation objectives. This compositional approach to simulation, called Functional Simulation (FSIL), supports the development of a variety of simulation strategies.

The Dynamic Simulation Environment (DSE) can be used to integrate legacy tools into a simulation, such integration techniques include; Gateway models which connect remote systems to a DSE simulation; and Model Compatibility Layers which import foreign models into a DSE simulation. Several Gateway and Model Compatibility Layers are available alongside the Model C Library, our C based implementation of the Dynamic Simulation Environment.
{{% /blocks/lead %}}


{{% blocks/section color="dark" type="row" %}}

{{% blocks/feature icon="fab fa-github" title="Getting Started!" url="docs/start/" %}}
Get started with the DSE.
{{% /blocks/feature %}}

{{% blocks/feature icon="fab fa-github" title="Contributions welcome!" url="https://github.com/boschglobal?q=dse&type=all" %}}
We do a <b>Pull Request</b> contributions workflow on **GitHub**.
{{% /blocks/feature %}}

{{% blocks/feature icon="fab fa-stack-overflow" title="Stack Overflow!" url="https://stackoverflow.com/questions/tagged/modelc" %}}
Ask technical questions.
{{% /blocks/feature %}}

{{% /blocks/section %}}
