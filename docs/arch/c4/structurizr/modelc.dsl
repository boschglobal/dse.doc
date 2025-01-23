group "ModelC Runtime" {
    codecApi = component "Codec API" "Binary-stream encode/decode."
    modelApi = component "Model API"
    signalApi = component "Signal API" "Vector access to signals."
}
codecLib = component "Network Codec Library" "Network codecs (PDU/CAN/ETH ...)."
modelApi -> signalApi "" "" RelApiCall
codecLib -> codecApi "Load MIMEtype defined codecs" "Shared Library" RelSharedLibrary
