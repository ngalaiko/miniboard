/**
 * @fileoverview gRPC-Web generated client stub for app.miniboard.codes.v1
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.app = {};
proto.app.miniboard = {};
proto.app.miniboard.codes = {};
proto.app.miniboard.codes.v1 = require('./codes_service_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.app.miniboard.codes.v1.CodesServiceClient =
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
proto.app.miniboard.codes.v1.CodesServicePromiseClient =
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
 *   !proto.app.miniboard.codes.v1.CreateCodeRequest,
 *   !proto.app.miniboard.codes.v1.Code>}
 */
const methodDescriptor_CodesService_CreateCode = new grpc.web.MethodDescriptor(
  '/app.miniboard.codes.v1.CodesService/CreateCode',
  grpc.web.MethodType.UNARY,
  proto.app.miniboard.codes.v1.CreateCodeRequest,
  proto.app.miniboard.codes.v1.Code,
  /**
   * @param {!proto.app.miniboard.codes.v1.CreateCodeRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.app.miniboard.codes.v1.Code.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.app.miniboard.codes.v1.CreateCodeRequest,
 *   !proto.app.miniboard.codes.v1.Code>}
 */
const methodInfo_CodesService_CreateCode = new grpc.web.AbstractClientBase.MethodInfo(
  proto.app.miniboard.codes.v1.Code,
  /**
   * @param {!proto.app.miniboard.codes.v1.CreateCodeRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.app.miniboard.codes.v1.Code.deserializeBinary
);


/**
 * @param {!proto.app.miniboard.codes.v1.CreateCodeRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.app.miniboard.codes.v1.Code)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.app.miniboard.codes.v1.Code>|undefined}
 *     The XHR Node Readable Stream
 */
proto.app.miniboard.codes.v1.CodesServiceClient.prototype.createCode =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/app.miniboard.codes.v1.CodesService/CreateCode',
      request,
      metadata || {},
      methodDescriptor_CodesService_CreateCode,
      callback);
};


/**
 * @param {!proto.app.miniboard.codes.v1.CreateCodeRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.app.miniboard.codes.v1.Code>}
 *     A native promise that resolves to the response
 */
proto.app.miniboard.codes.v1.CodesServicePromiseClient.prototype.createCode =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/app.miniboard.codes.v1.CodesService/CreateCode',
      request,
      metadata || {},
      methodDescriptor_CodesService_CreateCode);
};


module.exports = proto.app.miniboard.codes.v1;

