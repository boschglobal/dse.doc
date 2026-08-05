// Copyright 2026 Robert Bosch GmbH
//
// SPDX-License-Identifier: Apache-2.0

#include <stddef.h>
#include <dse/clib/util/yaml.h>
#include <dse/ncodec/codec.h>
#include <dse/ncodec/stream/stream.h>
#include <dse/pdunet/pdunet.h>

#define ARRAY_SIZE(x) (sizeof(x) / sizeof((x)[0]))
#define BUFFER_LEN    4096
#define MIMETYPE                                                               \
    "application/x-automotive-bus;"                                            \
    "interface=stream;type=pdu;schema=fbs;"                                    \
    "ecu_id=5;loopback=1"

int run_pdunet_example(const char* yaml)
{
    /* Define signals. */
    const char* signal_names[] = { "Vehicle.Speed", "Vehicle.Acceleration" };
    double      signal_values[] = { 0.0, 0.0 };

    /* Load Network YAML. */
    YamlDocList* dl = dse_yaml_load_file(NULL, yaml, NULL);
    void* doc = dse_yaml_find_doc_in_doclist(dl, "Network", NULL, NULL, 0);

    /* Create NCodec. */
    NSTREAM* stream = ncodec_buffer_stream_create(BUFFER_LEN);
    NCODEC*  nc = ncodec_open(MIMETYPE, stream);
    if (nc == NULL) return 1;

    /* Create and map PDUNet. */
    PDUNET* net = pdunet_create(nc, doc, 0.0005, NULL, NULL);
    if (net == NULL) {
        ncodec_close(nc);
        return 2;
    }
    int rc = pdunet_map_signals(
        net, "PDUNet", ARRAY_SIZE(signal_names), signal_names, signal_values);
    if (rc != 0) {
        pdunet_destroy(net);
        ncodec_close(nc);
        return 3;
    }

    /* Example simulation step. */
    double simulation_time = 0.0;
    signal_values[0] = 50.0;
    signal_values[1] = 1.2;
    pdunet_tx(net, NULL, NULL, NULL, simulation_time);
    pdunet_rx(net, NULL, NULL, NULL);

    /* Cleanup. */
    pdunet_destroy(net);
    ncodec_close(nc);
    return 0;
}
