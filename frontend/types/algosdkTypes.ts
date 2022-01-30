// @ts-nocheck
/**
 * For some reason, nextjs failed to import types to utils/api.ts from algosdk
 * so these are copied to this codebase to solve the problem
 */

/**
 * Base class for models
 */

/* eslint-disable no-underscore-dangle,camelcase */
function _is_primitive(val: any): val is string | boolean | number | bigint {
  /* eslint-enable no-underscore-dangle,camelcase */
  return (
    val === undefined ||
    val == null ||
    (typeof val !== "object" && typeof val !== "function")
  );
}

/* eslint-disable no-underscore-dangle,camelcase,no-redeclare,no-unused-vars */
function _get_obj_for_encoding(
  val: Function,
  binary: boolean
): Record<string, any>;
function _get_obj_for_encoding(val: any[], binary: boolean): any[];
function _get_obj_for_encoding(
  val: Record<string, any>,
  binary: boolean
): Record<string, any>;
function _get_obj_for_encoding(val: any, binary: boolean): any {
  /* eslint-enable no-underscore-dangle,camelcase,no-redeclare,no-unused-vars */
  let targetPropValue: any;

  if (val instanceof Uint8Array) {
    targetPropValue = binary ? val : Buffer.from(val).toString("base64");
  } else if (typeof val.get_obj_for_encoding === "function") {
    targetPropValue = val.get_obj_for_encoding(binary);
  } else if (Array.isArray(val)) {
    targetPropValue = [];
    for (const elem of val) {
      targetPropValue.push(_get_obj_for_encoding(elem, binary));
    }
  } else if (typeof val === "object") {
    const obj = {};
    for (const prop of Object.keys(val)) {
      obj[prop] = _get_obj_for_encoding(val[prop], binary);
    }
    targetPropValue = obj;
  } else if (_is_primitive(val)) {
    targetPropValue = val;
  } else {
    throw new Error(`Unsupported value: ${String(val)}`);
  }
  return targetPropValue;
}

export default class BaseModel {
  /* eslint-disable no-underscore-dangle,camelcase */
  attribute_map: Record<string, string>;

  /**
   * Get an object ready for encoding to either JSON or msgpack.
   * @param binary - Use true to indicate that the encoding can handle raw binary objects
   *   (Uint8Arrays). Use false to indicate that raw binary objects should be converted to base64
   *   strings. True should be used for objects that will be encoded with msgpack, and false should
   *   be used for objects that will be encoded with JSON.
   */
  get_obj_for_encoding(binary = false) {
    /* eslint-enable no-underscore-dangle,camelcase */
    const obj: Record<string, any> = {};

    for (const prop of Object.keys(this.attribute_map)) {
      const name = this.attribute_map[prop];
      const value = this[prop];

      if (typeof value !== "undefined") {
        obj[name] =
          value === null ? null : _get_obj_for_encoding(value, binary);
      }
    }

    return obj;
  }
}

/**
 * Application index and its parameters
 */
export class Application extends BaseModel {
  /**
   * (appidx) application index.
   */
  public id: number | bigint;

  /**
   * (appparams) application parameters.
   */
  public params: ApplicationParams;

  /**
   * Creates a new `Application` object.
   * @param id - (appidx) application index.
   * @param params - (appparams) application parameters.
   */
  constructor(id: number | bigint, params: ApplicationParams) {
    super();
    this.id = id;
    this.params = params;

    this.attribute_map = {
      id: "id",
      params: "params",
    };
  }
}

/**
 * Stores local state associated with an application.
 */
export class ApplicationLocalState extends BaseModel {
  /**
   * The application which this local state is for.
   */
  public id: number | bigint;

  /**
   * (hsch) schema.
   */
  public schema: ApplicationStateSchema;

  /**
   * (tkv) storage.
   */
  public keyValue?: TealKeyValue[];

  /**
   * Creates a new `ApplicationLocalState` object.
   * @param id - The application which this local state is for.
   * @param schema - (hsch) schema.
   * @param keyValue - (tkv) storage.
   */
  constructor(
    id: number | bigint,
    schema: ApplicationStateSchema,
    keyValue?: TealKeyValue[]
  ) {
    super();
    this.id = id;
    this.schema = schema;
    this.keyValue = keyValue;

    this.attribute_map = {
      id: "id",
      schema: "schema",
      keyValue: "key-value",
    };
  }
}

/**
 * Stores the global information associated with an application.
 */
export class ApplicationParams extends BaseModel {
  /**
   * (approv) approval program.
   */
  public approvalProgram: Uint8Array;

  /**
   * (clearp) approval program.
   */
  public clearStateProgram: Uint8Array;

  /**
   * The address that created this application. This is the address where the
   * parameters and global state for this application can be found.
   */
  public creator: string;

  /**
   * (epp) the amount of extra program pages available to this app.
   */
  public extraProgramPages?: number | bigint;

  /**
   * [\gs) global schema
   */
  public globalState?: TealKeyValue[];

  /**
   * [\gsch) global schema
   */
  public globalStateSchema?: ApplicationStateSchema;

  /**
   * [\lsch) local schema
   */
  public localStateSchema?: ApplicationStateSchema;

  /**
   * Creates a new `ApplicationParams` object.
   * @param approvalProgram - (approv) approval program.
   * @param clearStateProgram - (clearp) approval program.
   * @param creator - The address that created this application. This is the address where the
   * parameters and global state for this application can be found.
   * @param extraProgramPages - (epp) the amount of extra program pages available to this app.
   * @param globalState - [\gs) global schema
   * @param globalStateSchema - [\gsch) global schema
   * @param localStateSchema - [\lsch) local schema
   */
  constructor({
    approvalProgram,
    clearStateProgram,
    creator,
    extraProgramPages,
    globalState,
    globalStateSchema,
    localStateSchema,
  }: {
    approvalProgram: string | Uint8Array;
    clearStateProgram: string | Uint8Array;
    creator: string;
    extraProgramPages?: number | bigint;
    globalState?: TealKeyValue[];
    globalStateSchema?: ApplicationStateSchema;
    localStateSchema?: ApplicationStateSchema;
  }) {
    super();
    this.approvalProgram =
      typeof approvalProgram === "string"
        ? new Uint8Array(Buffer.from(approvalProgram, "base64"))
        : approvalProgram;
    this.clearStateProgram =
      typeof clearStateProgram === "string"
        ? new Uint8Array(Buffer.from(clearStateProgram, "base64"))
        : clearStateProgram;
    this.creator = creator;
    this.extraProgramPages = extraProgramPages;
    this.globalState = globalState;
    this.globalStateSchema = globalStateSchema;
    this.localStateSchema = localStateSchema;

    this.attribute_map = {
      approvalProgram: "approval-program",
      clearStateProgram: "clear-state-program",
      creator: "creator",
      extraProgramPages: "extra-program-pages",
      globalState: "global-state",
      globalStateSchema: "global-state-schema",
      localStateSchema: "local-state-schema",
    };
  }
}

/**
 * Specifies maximums on the number of each type that may be stored.
 */
export class ApplicationStateSchema extends BaseModel {
  /**
   * (nui) num of uints.
   */
  public numUint: number | bigint;

  /**
   * (nbs) num of byte slices.
   */
  public numByteSlice: number | bigint;

  /**
   * Creates a new `ApplicationStateSchema` object.
   * @param numUint - (nui) num of uints.
   * @param numByteSlice - (nbs) num of byte slices.
   */
  constructor(numUint: number | bigint, numByteSlice: number | bigint) {
    super();
    this.numUint = numUint;
    this.numByteSlice = numByteSlice;

    this.attribute_map = {
      numUint: "num-uint",
      numByteSlice: "num-byte-slice",
    };
  }
}

/**
 * Specifies both the unique identifier and the parameters for an asset
 */
export class Asset extends BaseModel {
  /**
   * unique asset identifier
   */
  public index: number | bigint;

  /**
   * AssetParams specifies the parameters for an asset.
   * (apar) when part of an AssetConfig transaction.
   * Definition:
   * data/transactions/asset.go : AssetParams
   */
  public params: AssetParams;

  /**
   * Creates a new `Asset` object.
   * @param index - unique asset identifier
   * @param params - AssetParams specifies the parameters for an asset.
   * (apar) when part of an AssetConfig transaction.
   * Definition:
   * data/transactions/asset.go : AssetParams
   */
  constructor(index: number | bigint, params: AssetParams) {
    super();
    this.index = index;
    this.params = params;

    this.attribute_map = {
      index: "index",
      params: "params",
    };
  }
}
