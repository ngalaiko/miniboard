/**
 * @fileoverview gRPC-Web generated client stub for app.miniboard.users.sources.v1
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.app = {};
proto.app.miniboard = {};
proto.app.miniboard.users = {};
proto.app.miniboard.users.sources = {};
proto.app.miniboard.users.sources.v1 = require('./sources_service_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.app.miniboard.users.sources.v1.SourcesServiceClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.app.miniboard.users.sources.v1.SourcesServicePromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.app.miniboard.users.sources.v1.CreateSourceRequest,
 *   !proto.app.miniboard.users.sources.v1.Source>}
 */
const methodDescriptor_SourcesService_CreateSource = new grpc.web.MethodDescriptor(
  '/app.miniboard.users.sources.v1.SourcesService/CreateSource',
  grpc.web.MethodType.UNARY,
  proto.app.miniboard.users.sources.v1.CreateSourceRequest,
  proto.app.miniboard.users.sources.v1.Source,
  /**
   * @param {!proto.app.miniboard.users.sources.v1.CreateSourceRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.app.miniboard.users.sources.v1.Source.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.app.miniboard.users.sources.v1.CreateSourceRequest,
 *   !proto.app.miniboard.users.sources.v1.Source>}
 */
const methodInfo_SourcesService_CreateSource = new grpc.web.AbstractClientBase.MethodInfo(
  proto.app.miniboard.users.sources.v1.Source,
  /**
   * @param {!proto.app.miniboard.users.sources.v1.CreateSourceRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.app.miniboard.users.sources.v1.Source.deserializeBinary
);


/**
 * @param {!proto.app.miniboard.users.sources.v1.CreateSourceRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.app.miniboard.users.sources.v1.Source)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.app.miniboard.users.sources.v1.Source>|undefined}
 *     The XHR Node Readable Stream
 */
proto.app.miniboard.users.sources.v1.SourcesServiceClient.prototype.createSource =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/app.miniboard.users.sources.v1.SourcesService/CreateSource',
      request,
      metadata || {},
      methodDescriptor_SourcesService_CreateSource,
      callback);
};


/**
 * @param {!proto.app.miniboard.users.sources.v1.CreateSourceRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.app.miniboard.users.sources.v1.Source>}
 *     A native promise that resolves to the response
 */
proto.app.miniboard.users.sources.v1.SourcesServicePromiseClient.prototype.createSource =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/app.miniboard.users.sources.v1.SourcesService/CreateSource',
      request,
      metadata || {},
      methodDescriptor_SourcesService_CreateSource);
};


module.exports = proto.app.miniboard.users.sources.v1;

