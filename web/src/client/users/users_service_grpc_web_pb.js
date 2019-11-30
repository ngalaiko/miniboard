/**
 * @fileoverview gRPC-Web generated client stub for app.miniboard.users.v1
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
proto.app.miniboard.users.v1 = require('./users_service_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.app.miniboard.users.v1.UsersServiceClient =
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
proto.app.miniboard.users.v1.UsersServicePromiseClient =
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
 *   !proto.app.miniboard.users.v1.GetUserRequest,
 *   !proto.app.miniboard.users.v1.User>}
 */
const methodDescriptor_UsersService_GetUser = new grpc.web.MethodDescriptor(
  '/app.miniboard.users.v1.UsersService/GetUser',
  grpc.web.MethodType.UNARY,
  proto.app.miniboard.users.v1.GetUserRequest,
  proto.app.miniboard.users.v1.User,
  /**
   * @param {!proto.app.miniboard.users.v1.GetUserRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.app.miniboard.users.v1.User.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.app.miniboard.users.v1.GetUserRequest,
 *   !proto.app.miniboard.users.v1.User>}
 */
const methodInfo_UsersService_GetUser = new grpc.web.AbstractClientBase.MethodInfo(
  proto.app.miniboard.users.v1.User,
  /**
   * @param {!proto.app.miniboard.users.v1.GetUserRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.app.miniboard.users.v1.User.deserializeBinary
);


/**
 * @param {!proto.app.miniboard.users.v1.GetUserRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.app.miniboard.users.v1.User)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.app.miniboard.users.v1.User>|undefined}
 *     The XHR Node Readable Stream
 */
proto.app.miniboard.users.v1.UsersServiceClient.prototype.getUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/app.miniboard.users.v1.UsersService/GetUser',
      request,
      metadata || {},
      methodDescriptor_UsersService_GetUser,
      callback);
};


/**
 * @param {!proto.app.miniboard.users.v1.GetUserRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.app.miniboard.users.v1.User>}
 *     A native promise that resolves to the response
 */
proto.app.miniboard.users.v1.UsersServicePromiseClient.prototype.getUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/app.miniboard.users.v1.UsersService/GetUser',
      request,
      metadata || {},
      methodDescriptor_UsersService_GetUser);
};


module.exports = proto.app.miniboard.users.v1;

