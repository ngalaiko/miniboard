/**
 * @fileoverview gRPC-Web generated client stub for app.miniboard.users.articles.v1
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');


var google_protobuf_empty_pb = require('google-protobuf/google/protobuf/empty_pb.js')

var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js')

var google_protobuf_field_mask_pb = require('google-protobuf/google/protobuf/field_mask_pb.js')

var google_protobuf_wrappers_pb = require('google-protobuf/google/protobuf/wrappers_pb.js')
const proto = {};
proto.app = {};
proto.app.miniboard = {};
proto.app.miniboard.users = {};
proto.app.miniboard.users.articles = {};
proto.app.miniboard.users.articles.v1 = require('./articles_service_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.app.miniboard.users.articles.v1.ArticlesServiceClient =
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
proto.app.miniboard.users.articles.v1.ArticlesServicePromiseClient =
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
 *   !proto.app.miniboard.users.articles.v1.ListArticlesRequest,
 *   !proto.app.miniboard.users.articles.v1.ListArticlesResponse>}
 */
const methodDescriptor_ArticlesService_ListArticles = new grpc.web.MethodDescriptor(
  '/app.miniboard.users.articles.v1.ArticlesService/ListArticles',
  grpc.web.MethodType.UNARY,
  proto.app.miniboard.users.articles.v1.ListArticlesRequest,
  proto.app.miniboard.users.articles.v1.ListArticlesResponse,
  /**
   * @param {!proto.app.miniboard.users.articles.v1.ListArticlesRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.app.miniboard.users.articles.v1.ListArticlesResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.app.miniboard.users.articles.v1.ListArticlesRequest,
 *   !proto.app.miniboard.users.articles.v1.ListArticlesResponse>}
 */
const methodInfo_ArticlesService_ListArticles = new grpc.web.AbstractClientBase.MethodInfo(
  proto.app.miniboard.users.articles.v1.ListArticlesResponse,
  /**
   * @param {!proto.app.miniboard.users.articles.v1.ListArticlesRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.app.miniboard.users.articles.v1.ListArticlesResponse.deserializeBinary
);


/**
 * @param {!proto.app.miniboard.users.articles.v1.ListArticlesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.app.miniboard.users.articles.v1.ListArticlesResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.app.miniboard.users.articles.v1.ListArticlesResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.app.miniboard.users.articles.v1.ArticlesServiceClient.prototype.listArticles =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/app.miniboard.users.articles.v1.ArticlesService/ListArticles',
      request,
      metadata || {},
      methodDescriptor_ArticlesService_ListArticles,
      callback);
};


/**
 * @param {!proto.app.miniboard.users.articles.v1.ListArticlesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.app.miniboard.users.articles.v1.ListArticlesResponse>}
 *     A native promise that resolves to the response
 */
proto.app.miniboard.users.articles.v1.ArticlesServicePromiseClient.prototype.listArticles =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/app.miniboard.users.articles.v1.ArticlesService/ListArticles',
      request,
      metadata || {},
      methodDescriptor_ArticlesService_ListArticles);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.app.miniboard.users.articles.v1.UpdateArticleRequest,
 *   !proto.app.miniboard.users.articles.v1.Article>}
 */
const methodDescriptor_ArticlesService_UpdateArticle = new grpc.web.MethodDescriptor(
  '/app.miniboard.users.articles.v1.ArticlesService/UpdateArticle',
  grpc.web.MethodType.UNARY,
  proto.app.miniboard.users.articles.v1.UpdateArticleRequest,
  proto.app.miniboard.users.articles.v1.Article,
  /**
   * @param {!proto.app.miniboard.users.articles.v1.UpdateArticleRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.app.miniboard.users.articles.v1.Article.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.app.miniboard.users.articles.v1.UpdateArticleRequest,
 *   !proto.app.miniboard.users.articles.v1.Article>}
 */
const methodInfo_ArticlesService_UpdateArticle = new grpc.web.AbstractClientBase.MethodInfo(
  proto.app.miniboard.users.articles.v1.Article,
  /**
   * @param {!proto.app.miniboard.users.articles.v1.UpdateArticleRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.app.miniboard.users.articles.v1.Article.deserializeBinary
);


/**
 * @param {!proto.app.miniboard.users.articles.v1.UpdateArticleRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.app.miniboard.users.articles.v1.Article)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.app.miniboard.users.articles.v1.Article>|undefined}
 *     The XHR Node Readable Stream
 */
proto.app.miniboard.users.articles.v1.ArticlesServiceClient.prototype.updateArticle =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/app.miniboard.users.articles.v1.ArticlesService/UpdateArticle',
      request,
      metadata || {},
      methodDescriptor_ArticlesService_UpdateArticle,
      callback);
};


/**
 * @param {!proto.app.miniboard.users.articles.v1.UpdateArticleRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.app.miniboard.users.articles.v1.Article>}
 *     A native promise that resolves to the response
 */
proto.app.miniboard.users.articles.v1.ArticlesServicePromiseClient.prototype.updateArticle =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/app.miniboard.users.articles.v1.ArticlesService/UpdateArticle',
      request,
      metadata || {},
      methodDescriptor_ArticlesService_UpdateArticle);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.app.miniboard.users.articles.v1.GetArticleRequest,
 *   !proto.app.miniboard.users.articles.v1.Article>}
 */
const methodDescriptor_ArticlesService_GetArticle = new grpc.web.MethodDescriptor(
  '/app.miniboard.users.articles.v1.ArticlesService/GetArticle',
  grpc.web.MethodType.UNARY,
  proto.app.miniboard.users.articles.v1.GetArticleRequest,
  proto.app.miniboard.users.articles.v1.Article,
  /**
   * @param {!proto.app.miniboard.users.articles.v1.GetArticleRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.app.miniboard.users.articles.v1.Article.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.app.miniboard.users.articles.v1.GetArticleRequest,
 *   !proto.app.miniboard.users.articles.v1.Article>}
 */
const methodInfo_ArticlesService_GetArticle = new grpc.web.AbstractClientBase.MethodInfo(
  proto.app.miniboard.users.articles.v1.Article,
  /**
   * @param {!proto.app.miniboard.users.articles.v1.GetArticleRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.app.miniboard.users.articles.v1.Article.deserializeBinary
);


/**
 * @param {!proto.app.miniboard.users.articles.v1.GetArticleRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.app.miniboard.users.articles.v1.Article)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.app.miniboard.users.articles.v1.Article>|undefined}
 *     The XHR Node Readable Stream
 */
proto.app.miniboard.users.articles.v1.ArticlesServiceClient.prototype.getArticle =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/app.miniboard.users.articles.v1.ArticlesService/GetArticle',
      request,
      metadata || {},
      methodDescriptor_ArticlesService_GetArticle,
      callback);
};


/**
 * @param {!proto.app.miniboard.users.articles.v1.GetArticleRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.app.miniboard.users.articles.v1.Article>}
 *     A native promise that resolves to the response
 */
proto.app.miniboard.users.articles.v1.ArticlesServicePromiseClient.prototype.getArticle =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/app.miniboard.users.articles.v1.ArticlesService/GetArticle',
      request,
      metadata || {},
      methodDescriptor_ArticlesService_GetArticle);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.app.miniboard.users.articles.v1.DeleteArticleRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_ArticlesService_DeleteArticle = new grpc.web.MethodDescriptor(
  '/app.miniboard.users.articles.v1.ArticlesService/DeleteArticle',
  grpc.web.MethodType.UNARY,
  proto.app.miniboard.users.articles.v1.DeleteArticleRequest,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.app.miniboard.users.articles.v1.DeleteArticleRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.app.miniboard.users.articles.v1.DeleteArticleRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_ArticlesService_DeleteArticle = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.app.miniboard.users.articles.v1.DeleteArticleRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.app.miniboard.users.articles.v1.DeleteArticleRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.app.miniboard.users.articles.v1.ArticlesServiceClient.prototype.deleteArticle =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/app.miniboard.users.articles.v1.ArticlesService/DeleteArticle',
      request,
      metadata || {},
      methodDescriptor_ArticlesService_DeleteArticle,
      callback);
};


/**
 * @param {!proto.app.miniboard.users.articles.v1.DeleteArticleRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     A native promise that resolves to the response
 */
proto.app.miniboard.users.articles.v1.ArticlesServicePromiseClient.prototype.deleteArticle =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/app.miniboard.users.articles.v1.ArticlesService/DeleteArticle',
      request,
      metadata || {},
      methodDescriptor_ArticlesService_DeleteArticle);
};


module.exports = proto.app.miniboard.users.articles.v1;

