/**
 * @fileoverview gRPC-Web generated client stub for app.miniboard.tokens.v1
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.app = {};
proto.app.miniboard = {};
proto.app.miniboard.tokens = {};
proto.app.miniboard.tokens.v1 = require('./tokens_service_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.app.miniboard.tokens.v1.TokensServiceClient =
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
proto.app.miniboard.tokens.v1.TokensServicePromiseClient =
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
 *   !proto.app.miniboard.tokens.v1.CreateTokenRequest,
 *   !proto.app.miniboard.tokens.v1.Token>}
 */
const methodDescriptor_TokensService_CreateToken = new grpc.web.MethodDescriptor(
  '/app.miniboard.tokens.v1.TokensService/CreateToken',
  grpc.web.MethodType.UNARY,
  proto.app.miniboard.tokens.v1.CreateTokenRequest,
  proto.app.miniboard.tokens.v1.Token,
  /**
   * @param {!proto.app.miniboard.tokens.v1.CreateTokenRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.app.miniboard.tokens.v1.Token.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.app.miniboard.tokens.v1.CreateTokenRequest,
 *   !proto.app.miniboard.tokens.v1.Token>}
 */
const methodInfo_TokensService_CreateToken = new grpc.web.AbstractClientBase.MethodInfo(
  proto.app.miniboard.tokens.v1.Token,
  /**
   * @param {!proto.app.miniboard.tokens.v1.CreateTokenRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.app.miniboard.tokens.v1.Token.deserializeBinary
);


/**
 * @param {!proto.app.miniboard.tokens.v1.CreateTokenRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.app.miniboard.tokens.v1.Token)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.app.miniboard.tokens.v1.Token>|undefined}
 *     The XHR Node Readable Stream
 */
proto.app.miniboard.tokens.v1.TokensServiceClient.prototype.createToken =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/app.miniboard.tokens.v1.TokensService/CreateToken',
      request,
      metadata || {},
      methodDescriptor_TokensService_CreateToken,
      callback);
};


/**
 * @param {!proto.app.miniboard.tokens.v1.CreateTokenRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.app.miniboard.tokens.v1.Token>}
 *     A native promise that resolves to the response
 */
proto.app.miniboard.tokens.v1.TokensServicePromiseClient.prototype.createToken =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/app.miniboard.tokens.v1.TokensService/CreateToken',
      request,
      metadata || {},
      methodDescriptor_TokensService_CreateToken);
};


module.exports = proto.app.miniboard.tokens.v1;

