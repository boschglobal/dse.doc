# Copyright 2025 Robert Bosch GmbH
#
# SPDX-License-Identifier: Apache-2.0

model {
    properties {
        "structurizr.groupSeparator" "/"
    }

    u = person "Simulation Engineer"
    fsil = softwareSystem "${C_MONIKER} Simulation" {
        tags "fsil" "FsilSystem"
        description "CoSim: fixed/variable step-size, events\nSignals: scalar, binary-streams\nRuntime: distributed multi-platform/multi-OS"

        simbus = container "SimBus" {
            tags "SimBus"
        }


        # Containers: Native
        # ==================
        model = container "Model" {
            !include modelc.dsl
            tags "native"
            description "Runs in single OS process."

            model = component "Model"
            modelApi -> model "Load and Step the model" "Shared Library" RelSharedLibrary
            signalApi -> model "Exchange scalar/binary signals" "scalar/binary vector" RelSignalExchange
            codecApi -> model "Encode/decode binary streams" "Flatbuffer/MsgPack" RelEncodedBinary
        }
        composite = container "Composite Model" {
            !include modelc.dsl
            tags "native"
            description "Loads several models into single OS process. Operation alogorithm: co-sim or composite/ordered (output feeds to input)."

            group "Model Stack" {
                controller = component "Controller"
                modelA = component "Model A"
                modelB = component "Model B"
            }
            modelApi -> modelA "Load the model" "Shared Library"
            modelApi -> modelB "Load the model" "Shared Library"
            signalApi -> controller "Exchange scalar/binary signals" "scalar/binary vector"
            codecApi -> modelA "Encode/decode binary streams" "Flatbuffer/MsgPack"
            codecApi -> modelB "Encode/decode binary streams" "Flatbuffer/MsgPack"
            controller -> modelA "Step"
            modelA -> modelB "Step"
            modelA -> controller "Exchange scalar/binary signals"
            modelB -> modelA "Exchange scalar/binary signals"
        }
        simbus -> model.modelApi  "Signal Exchange: Scalar/Binary" "Flatbuffer/MsgPack" RelSimBusConnection
        simbus -> composite.modelApi  "Signal Exchange: Scalar/Binary" "Flatbuffer/MsgPack" RelSimBusConnection


        # Containers: Virtual ECU
        # =======================
        bootloader = container "Bootloader Model" {
            !include modelc.dsl
            tags "vecu"
            description "Loads a Virtual ECU using a Bootloader."

            softecu = component "Virtual ECU" "SoftECU based Virtual ECU."
            bootloader = component "Bollogg" "SoftECU bootloader."
            bootloader -> softecu "Load and operate the SoftECU" "Shared Library/C89" RelSharedLibrary
            modelApi -> bootloader "Load the bootloader" "Shared Library" RelSharedLibrary
            signalApi -> bootloader "Exchange scalar/binary signals" "scalar/binary vector" RelSignalExchange
            codecApi -> bootloader "Encode/decode binary streams" "Flatbuffer/MsgPack" RelEncodedBinary
        }
        runnable = container "Runnable Model" {
            !include modelc.dsl
            tags "vecu"
            description "Task Scheduler with in-direct Signal/PDU/Network interfaces."

            task = component "Virtual ECU" "Externally scheduled Task based Virtual ECU."
            group "Runnable" {
                marshaler = component "Marshaler" "Low level data mapping to C Structures."
                scheduler = component "Scheduler" "Tick based schedule, multiple of CoSim step size."
            }
            signalApi -> marshaler "Exchange scalar/binary signals" "scalar/binary vector" RelSignalExchange
            modelApi -> scheduler "Load and Step the model" "Shared Library" RelSharedLibrary
            codecApi -> marshaler "Encode/decode binary streams" "Flatbuffer/MsgPack" RelEncodedBinary
            marshaler -> task "Access memory space" "Structure mapping" RelSharedLibrary
            scheduler -> task "Execute tasks according to schedule" "Function call" RelSharedLibrary
        }
        simbus -> bootloader.modelApi "Signal Exchange: Scalar/Binary" "Flatbuffer/MsgPack" RelSimBusConnection
        simbus -> runnable.modelApi "Signal Exchange: Scalar/Binary" "Flatbuffer/MsgPack" RelSimBusConnection


        # Containers: Interop
        # ===================
        gateway = container "Gateway Model" {
            tags "interop" "remote"
            description "Runs in the external/remote simulation and provides a logical point-of-presence"
        }
        mcl = container "MCL Model" {
            tags "interop" "external"
            description "Adapter for CoSim interface"
        }
        group "FMI Interop" {
            fmimcl = container "FMI MCL Model" {
                !include modelc.dsl
                tags "interop" "external"
                description "Adapter to FMU CoSim interface\n(FMI2/FMI3)."

                fmu = component "FMU" "CoSim FMU (FMI2/FMI3)."
                group "FMI MCL Model" {
                    adapter = component "Adapter" "Adapter for operation of FMU via versioned FMI API.\nVariable exchange."
                    fmimcl = component "FMI MCL" "Model compatability library for FMI Standard FMUs."
                }

                adapter -> fmu "Load the FMU" "Shared Library" RelSharedLibrary
                fmimcl -> adapter "Select adapter based on FMU version" "" RelSharedLibrary
                modelApi -> fmimcl "Load the MCL" "Shared Library" RelSharedLibrary
                signalApi -> fmimcl "Exchange scalar/binary signals" "scalar/binary vector" RelSignalExchange
                codecApi -> fmimcl "Encode/decode binary streams" "Flatbuffer/MsgPack" RelEncodedBinary
            }
            fmigateway = container "FMI Gateway FMU" {
                !include modelc.dsl
                tags "interop" "remote"
                description "Runs in the external/remote simulation."

                gateway = component "Gateway" ""
                modelApi -> gateway "Connect with ${C_MONIKER} simulation" "shared/static library" RelSharedLibrary
                signalApi -> gateway "Exchange scalar/binary signals" "scalar/binary vector" RelSignalExchange
                codecApi -> gateway "Encode/decode binary streams" "Flatbuffer/MsgPack" RelEncodedBinary
            }
        }
        simbus -> gateway  "Signal Exchange: Scalar/Binary" "Flatbuffer/MsgPack" RelSimBusConnection
        simbus -> mcl  "Signal Exchange: Scalar/Binary" "Flatbuffer/MsgPack" RelSimBusConnection
        simbus -> fmigateway.modelApi  "Signal Exchange: Scalar/Binary" "Flatbuffer/MsgPack" RelSimBusConnection
        simbus -> fmimcl.modelApi  "Signal Exchange: Scalar/Binary" "Flatbuffer/MsgPack" RelSimBusConnection

    }
    fsilModel = element "${C_MONIKER} Model" "Model" "Vector based signal access.\nBinary streams (Codecs/MIMEtype)" "DseModel"


    # System Context - External elements (project deliverable, systems).
    # ==================================
    group "Virtual ECU" {
        softEcu = element "Virtual ECU (SoftECU)" "VECU" "SoftECU based Virtual ECU\n(Legacy Autosar: ASW, CSW ...)" "VirtualECU"
        soaEcu = element "Virtual ECU (SOA)" "VECU" "SOA based Virtual ECU\n(Adaptive Autosar: POSIX, Task ...)" "VirtualECU"
    }
    group "Interop" {
        remoteSim = softwareSystem "External Simulation" "External Simulation Systems\nRemote Simulations\n(Simulink, CarMaker, ECU-Test ...)" "ExternalSystem"
        externalModel = element "External Model" "Model" "Foreign model standards\n(Simulink, FMI ...)" "ExternalModel"
        fmuModel = element "FMU" "FMI2/FMI3" "FMU Package" ExternalModel,FMI
        fmiSimulation = softwareSystem "FMI Simulation" "FMI2/FMI3" ExternalSystem,FMI
    }

    # Relationships to/from external elements.
    fsil.bootloader -> softEcu "Runs" "Shared Library"
    fsil.runnable -> soaEcu "Runs" "Shared Library"
    fsil.gateway -> remoteSim "Gateway Connection" "TCP/Redis/Flatbuffer"
    fsil.mcl -> externalModel "Load via Model Compatability Library" "Shared Library"
    fsil.fmimcl -> fmuModel "Loads" "Shared Library"
    fsil.fmigateway -> fmiSimulation "Gateway Connection" "TCP/Redis/Flatbuffer"

    # Relationships to/from FSIL elements.
    fsil.model -> fsilModel "Loads" "Shared Library"
    fsil.composite -> fsilModel "Loads" "Shared Library"


    u -> fsil "Operates"
}

views {
    systemContext fsil "SystemContext" {
        include *
        exclude element.tag==FMI
        autolayout
        description "${C_MONIKER} System Context.
    }

    container fsil "Container-Native" {
        include ->element.tag==native->
        autolayout
    }
    container fsil "Container-VirtualEcu" {
        include ->element.tag==vecu->
        autolayout
    }
    container fsil "Container-Interop" {
        include ->element.tag==interop->
        autolayout
    }

    component fsil.model "Component-Model" {
        include *
        autolayout bt
    }
    component fsil.composite "Component-Composite-Model" {
        include *
        autolayout bt
    }
    component fsil.bootloader "Component-Bootloader" {
        include *
        autolayout bt
    }
    component fsil.runnable "Component-Runnable" {
        include *
        autolayout bt
    }
    component fsil.fmimcl "Component-FMI-MCL" {
        include *
        autolayout bt
    }
    component fsil.fmigateway "Component-FMI-Gateway" {
        include *
        autolayout tb
    }


    styles {
        element "Element" {
            shape RoundedBox
        }
        element "Person" {
            background #08427b
            color #ffffff
            shape person
        }
        element "Software System" {
            background #1168bd
            color #ffffff
        }
        element "Container" {
            background #438dd5
            color #ffffff
        }


        # Classification tags (for Legend)
        element "VirtualECU" {
            background #999999
        }
        element "SimBus" {
            background #1168bd
            color #ffffff
            opacity 50
            height 150
        }
        relationship "RelSimBusConnection" {
            color #1168bd
            opacity 50
            style solid
            routing Orthogonal
        }
        element "DseModel" {
            background #1168bd
            color #ffffff
        }
        element "ExternalModel" {
            background #AAAAAA
        }
        element "FsilSystem" {
            width 800
        }
        element "ExternalSystem" {
            background #BBBBBB
            color #000000
        }

        relationship "RelApiCall" {
            #color #1168bd
            style solid
        }
        relationship "RelEncodedBinary" {
            #color #1168bd
            style solid
        }
        relationship "RelSharedLibrary" {
            #color #1168bd
            style solid
        }
        relationship "RelSignalExchange" {
            #color #1168bd
            style solid
        }


    }
}

configuration {
    scope softwaresystem
}
