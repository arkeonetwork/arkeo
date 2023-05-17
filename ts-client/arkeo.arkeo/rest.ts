/* eslint-disable */
/* tslint:disable */
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface ArkeoContract {
  /** @format byte */
  provider?: string;

  /** @format int32 */
  service?: number;

  /** @format byte */
  client?: string;

  /** @format byte */
  delegate?: string;
  type?: ArkeoContractType;

  /** @format int64 */
  height?: string;

  /** @format int64 */
  duration?: string;

  /**
   * Coin defines a token with a denomination and an amount.
   *
   * NOTE: The amount field is an Int which implements the custom method
   * signatures required by gogoproto.
   */
  rate?: V1Beta1Coin;
  deposit?: string;
  paid?: string;

  /** @format int64 */
  nonce?: string;

  /** @format int64 */
  settlement_height?: string;

  /** @format uint64 */
  id?: string;

  /** @format int64 */
  settlement_duration?: string;
  authorization?: ArkeoContractAuthorization;

  /** @format int64 */
  queries_per_minute?: string;
}

export enum ArkeoContractAuthorization {
  STRICT = "STRICT",
  OPEN = "OPEN",
}

export enum ArkeoContractType {
  SUBSCRIPTION = "SUBSCRIPTION",
  PAY_AS_YOU_GO = "PAY_AS_YOU_GO",
}

export type ArkeoMsgBondProviderResponse = object;

export type ArkeoMsgClaimContractIncomeResponse = object;

export type ArkeoMsgCloseContractResponse = object;

export type ArkeoMsgModProviderResponse = object;

export type ArkeoMsgOpenContractResponse = object;

export type ArkeoMsgSetVersionResponse = object;

/**
 * Params defines the parameters for the module.
 */
export type ArkeoParams = object;

export interface ArkeoProvider {
  /** @format byte */
  pub_key?: string;

  /** @format int32 */
  service?: number;
  metadata_uri?: string;

  /** @format uint64 */
  metadata_nonce?: string;
  status?: ArkeoProviderStatus;

  /** @format int64 */
  min_contract_duration?: string;

  /** @format int64 */
  max_contract_duration?: string;
  subscription_rate?: V1Beta1Coin[];
  pay_as_you_go_rate?: V1Beta1Coin[];
  bond?: string;

  /** @format int64 */
  last_update?: string;

  /** @format int64 */
  settlement_duration?: string;
}

export enum ArkeoProviderStatus {
  OFFLINE = "OFFLINE",
  ONLINE = "ONLINE",
}

export interface ArkeoQueryActiveContractResponse {
  contract?: ArkeoContract;
}

export interface ArkeoQueryAllContractResponse {
  contract?: ArkeoContract[];

  /**
   * PageResponse is to be embedded in gRPC response messages where the
   * corresponding request message has used PageRequest.
   *
   *  message SomeResponse {
   *          repeated Bar results = 1;
   *          PageResponse page = 2;
   *  }
   */
  pagination?: V1Beta1PageResponse;
}

export interface ArkeoQueryAllProviderResponse {
  provider?: ArkeoProvider[];

  /**
   * PageResponse is to be embedded in gRPC response messages where the
   * corresponding request message has used PageRequest.
   *
   *  message SomeResponse {
   *          repeated Bar results = 1;
   *          PageResponse page = 2;
   *  }
   */
  pagination?: V1Beta1PageResponse;
}

export interface ArkeoQueryFetchContractResponse {
  contract?: ArkeoContract;
}

export interface ArkeoQueryFetchProviderResponse {
  provider?: ArkeoProvider;
}

/**
 * QueryParamsResponse is response type for the Query/Params RPC method.
 */
export interface ArkeoQueryParamsResponse {
  /** params holds all the parameters of this module. */
  params?: ArkeoParams;
}

export interface ProtobufAny {
  "@type"?: string;
}

export interface RpcStatus {
  /** @format int32 */
  code?: number;
  message?: string;
  details?: ProtobufAny[];
}

/**
* Coin defines a token with a denomination and an amount.

NOTE: The amount field is an Int which implements the custom method
signatures required by gogoproto.
*/
export interface V1Beta1Coin {
  denom?: string;
  amount?: string;
}

/**
* message SomeRequest {
         Foo some_parameter = 1;
         PageRequest pagination = 2;
 }
*/
export interface V1Beta1PageRequest {
  /**
   * key is a value returned in PageResponse.next_key to begin
   * querying the next page most efficiently. Only one of offset or key
   * should be set.
   * @format byte
   */
  key?: string;

  /**
   * offset is a numeric offset that can be used when key is unavailable.
   * It is less efficient than using key. Only one of offset or key should
   * be set.
   * @format uint64
   */
  offset?: string;

  /**
   * limit is the total number of results to be returned in the result page.
   * If left empty it will default to a value to be set by each app.
   * @format uint64
   */
  limit?: string;

  /**
   * count_total is set to true  to indicate that the result set should include
   * a count of the total number of items available for pagination in UIs.
   * count_total is only respected when offset is used. It is ignored when key
   * is set.
   */
  count_total?: boolean;

  /**
   * reverse is set to true if results are to be returned in the descending order.
   *
   * Since: cosmos-sdk 0.43
   */
  reverse?: boolean;
}

/**
* PageResponse is to be embedded in gRPC response messages where the
corresponding request message has used PageRequest.

 message SomeResponse {
         repeated Bar results = 1;
         PageResponse page = 2;
 }
*/
export interface V1Beta1PageResponse {
  /**
   * next_key is the key to be passed to PageRequest.key to
   * query the next page most efficiently. It will be empty if
   * there are no more results.
   * @format byte
   */
  next_key?: string;

  /**
   * total is total number of results available if PageRequest.count_total
   * was set, its value is undefined otherwise
   * @format uint64
   */
  total?: string;
}

import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse, ResponseType } from "axios";

export type QueryParamsType = Record<string | number, any>;

export interface FullRequestParams extends Omit<AxiosRequestConfig, "data" | "params" | "url" | "responseType"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseType;
  /** request body */
  body?: unknown;
}

export type RequestParams = Omit<FullRequestParams, "body" | "method" | "query" | "path">;

export interface ApiConfig<SecurityDataType = unknown> extends Omit<AxiosRequestConfig, "data" | "cancelToken"> {
  securityWorker?: (
    securityData: SecurityDataType | null,
  ) => Promise<AxiosRequestConfig | void> | AxiosRequestConfig | void;
  secure?: boolean;
  format?: ResponseType;
}

export enum ContentType {
  Json = "application/json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
}

export class HttpClient<SecurityDataType = unknown> {
  public instance: AxiosInstance;
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private secure?: boolean;
  private format?: ResponseType;

  constructor({ securityWorker, secure, format, ...axiosConfig }: ApiConfig<SecurityDataType> = {}) {
    this.instance = axios.create({ ...axiosConfig, baseURL: axiosConfig.baseURL || "" });
    this.secure = secure;
    this.format = format;
    this.securityWorker = securityWorker;
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  private mergeRequestParams(params1: AxiosRequestConfig, params2?: AxiosRequestConfig): AxiosRequestConfig {
    return {
      ...this.instance.defaults,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...(this.instance.defaults.headers || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  private createFormData(input: Record<string, unknown>): FormData {
    return Object.keys(input || {}).reduce((formData, key) => {
      const property = input[key];
      formData.append(
        key,
        property instanceof Blob
          ? property
          : typeof property === "object" && property !== null
          ? JSON.stringify(property)
          : `${property}`,
      );
      return formData;
    }, new FormData());
  }

  public request = async <T = any, _E = any>({
    secure,
    path,
    type,
    query,
    format,
    body,
    ...params
  }: FullRequestParams): Promise<AxiosResponse<T>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const responseFormat = (format && this.format) || void 0;

    if (type === ContentType.FormData && body && body !== null && typeof body === "object") {
      requestParams.headers.common = { Accept: "*/*" };
      requestParams.headers.post = {};
      requestParams.headers.put = {};

      body = this.createFormData(body as Record<string, unknown>);
    }

    return this.instance.request({
      ...requestParams,
      headers: {
        ...(type && type !== ContentType.FormData ? { "Content-Type": type } : {}),
        ...(requestParams.headers || {}),
      },
      params: query,
      responseType: responseFormat,
      data: body,
      url: path,
    });
  };
}

/**
 * @title arkeo/arkeo/events.proto
 * @version version not set
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  /**
   * No description
   *
   * @tags Query
   * @name QueryActiveContract
   * @summary Queries an active contract by spender, provider and service.
   * @request GET:/arkeo/active-contract/{provider}/{service}/{spender}
   */
  queryActiveContract = (provider: string, service: string, spender: string, params: RequestParams = {}) =>
    this.request<ArkeoQueryActiveContractResponse, RpcStatus>({
      path: `/arkeo/active-contract/${provider}/${service}/${spender}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryFetchContract
   * @request GET:/arkeo/contract/{contract_id}
   */
  queryFetchContract = (contractId: string, params: RequestParams = {}) =>
    this.request<ArkeoQueryFetchContractResponse, RpcStatus>({
      path: `/arkeo/contract/${contractId}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryContractAll
   * @request GET:/arkeo/contracts
   */
  queryContractAll = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<ArkeoQueryAllContractResponse, RpcStatus>({
      path: `/arkeo/contracts`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryParams
   * @summary Parameters queries the parameters of the module.
   * @request GET:/arkeo/params
   */
  queryParams = (params: RequestParams = {}) =>
    this.request<ArkeoQueryParamsResponse, RpcStatus>({
      path: `/arkeo/params`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryFetchProvider
   * @request GET:/arkeo/provider/{pubkey}/{service}
   */
  queryFetchProvider = (pubkey: string, service: string, params: RequestParams = {}) =>
    this.request<ArkeoQueryFetchProviderResponse, RpcStatus>({
      path: `/arkeo/provider/${pubkey}/${service}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryProviderAll
   * @request GET:/arkeo/providers
   */
  queryProviderAll = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<ArkeoQueryAllProviderResponse, RpcStatus>({
      path: `/arkeo/providers`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });
}
