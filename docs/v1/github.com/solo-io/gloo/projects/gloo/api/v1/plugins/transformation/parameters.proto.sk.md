<!-- Code generated by solo-kit. DO NOT EDIT. -->

### Package: `transformation.plugins.gloo.solo.io` 
##### Types:


- [Parameters](#Parameters)
  



##### Source File: [github.com/solo-io/gloo/projects/gloo/api/v1/plugins/transformation/parameters.proto](https://github.com/solo-io/gloo/blob/master/projects/gloo/api/v1/plugins/transformation/parameters.proto)





---
### <a name="Parameters">Parameters</a>



```yaml
"headers": map<string, string>
"path": .google.protobuf.StringValue

```

| Field | Type | Description | Default |
| ----- | ---- | ----------- |----------- | 
| `headers` | `map<string, string>` | headers that will be used to extract data for processing output templates Gloo will search for parameters by their name in header value strings, enclosed in single curly braces Example: extensions: parameters: headers: x-user-id: { userId } |  |
| `path` | [.google.protobuf.StringValue](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/string-value) | part of the (or the entire) path that will be used extract data for processing output templates Gloo will search for parameters by their name in header value strings, enclosed in single curly braces Example: extensions: parameters: path: /users/{ userId } |  |





<!-- Start of HubSpot Embed Code -->
<script type="text/javascript" id="hs-script-loader" async defer src="//js.hs-scripts.com/5130874.js"></script>
<!-- End of HubSpot Embed Code -->
