export namespace biz {
	
	export class AiChannelQuery {
	    page?: model.PageInfo;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new AiChannelQuery(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.page = this.convertValues(source["page"], model.PageInfo);
	        this.name = source["name"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DbDataSourceQuery {
	    page?: model.PageInfo;
	    name: string;
	    dbname: string;
	
	    static createFrom(source: any = {}) {
	        return new DbDataSourceQuery(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.page = this.convertValues(source["page"], model.PageInfo);
	        this.name = source["name"];
	        this.dbname = source["dbname"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DocumentQuery {
	    page?: model.PageInfo;
	    collectionName: string;
	    filename: string;
	    content: string;
	    metadata: string;
	
	    static createFrom(source: any = {}) {
	        return new DocumentQuery(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.page = this.convertValues(source["page"], model.PageInfo);
	        this.collectionName = source["collectionName"];
	        this.filename = source["filename"];
	        this.content = source["content"];
	        this.metadata = source["metadata"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FunctionImplQuery {
	    page?: model.PageInfo;
	    name: string;
	    title: string;
	
	    static createFrom(source: any = {}) {
	        return new FunctionImplQuery(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.page = this.convertValues(source["page"], model.PageInfo);
	        this.name = source["name"];
	        this.title = source["title"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FunctionQuery {
	    page?: model.PageInfo;
	    funcName: string;
	    roleId: string;
	
	    static createFrom(source: any = {}) {
	        return new FunctionQuery(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.page = this.convertValues(source["page"], model.PageInfo);
	        this.funcName = source["funcName"];
	        this.roleId = source["roleId"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FwConfigQuery {
	    page?: model.PageInfo;
	    name: string;
	    flow_type: string;
	    tenant_id: string;
	
	    static createFrom(source: any = {}) {
	        return new FwConfigQuery(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.page = this.convertValues(source["page"], model.PageInfo);
	        this.name = source["name"];
	        this.flow_type = source["flow_type"];
	        this.tenant_id = source["tenant_id"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PromptQuery {
	    page?: model.PageInfo;
	    roleName: string;
	
	    static createFrom(source: any = {}) {
	        return new PromptQuery(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.page = this.convertValues(source["page"], model.PageInfo);
	        this.roleName = source["roleName"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace dv {
	
	export class TableInfo {
	    tableName: string;
	    tableDesc: string;
	
	    static createFrom(source: any = {}) {
	        return new TableInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tableName = source["tableName"];
	        this.tableDesc = source["tableDesc"];
	    }
	}

}

export namespace entity {
	
	export class AiChannel {
	    id: string;
	    name: string;
	    url: string;
	    token: string;
	    maxToken: number;
	    models: string;
	    priority: number;
	    remark: string;
	    originalUrl: string;
	
	    static createFrom(source: any = {}) {
	        return new AiChannel(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.url = source["url"];
	        this.token = source["token"];
	        this.maxToken = source["maxToken"];
	        this.models = source["models"];
	        this.priority = source["priority"];
	        this.remark = source["remark"];
	        this.originalUrl = source["originalUrl"];
	    }
	}
	export class Config {
	    id: string;
	    content: string;
	    legalStatement: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.content = source["content"];
	        this.legalStatement = source["legalStatement"];
	    }
	}
	export class DbConfig {
	    id: string;
	    dbType: string;
	    dbname: string;
	    url: string;
	    username: string;
	    password: string;
	
	    static createFrom(source: any = {}) {
	        return new DbConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.dbType = source["dbType"];
	        this.dbname = source["dbname"];
	        this.url = source["url"];
	        this.username = source["username"];
	        this.password = source["password"];
	    }
	}
	export class Document {
	    id: string;
	    collectionName: string;
	    filename: string;
	    content: string;
	    metadata: string;
	
	    static createFrom(source: any = {}) {
	        return new Document(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.collectionName = source["collectionName"];
	        this.filename = source["filename"];
	        this.content = source["content"];
	        this.metadata = source["metadata"];
	    }
	}
	export class Function {
	    id: string;
	    funcName: string;
	    funcContent: string;
	    roleId: string;
	
	    static createFrom(source: any = {}) {
	        return new Function(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.funcName = source["funcName"];
	        this.funcContent = source["funcContent"];
	        this.roleId = source["roleId"];
	    }
	}
	export class FunctionImpl {
	    id: string;
	    title: string;
	    name: string;
	    url: string;
	    method: string;
	    contentType: string;
	    headers: string;
	
	    static createFrom(source: any = {}) {
	        return new FunctionImpl(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.name = source["name"];
	        this.url = source["url"];
	        this.method = source["method"];
	        this.contentType = source["contentType"];
	        this.headers = source["headers"];
	    }
	}
	export class FwConfig {
	    id: string;
	    name: string;
	    remark: string;
	    content: string;
	    flowType: string;
	    flowState: string;
	
	    static createFrom(source: any = {}) {
	        return new FwConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.remark = source["remark"];
	        this.content = source["content"];
	        this.flowType = source["flowType"];
	        this.flowState = source["flowState"];
	    }
	}
	export class Prompt {
	    id: string;
	    roleName: string;
	    prompt: string;
	
	    static createFrom(source: any = {}) {
	        return new Prompt(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.roleName = source["roleName"];
	        this.prompt = source["prompt"];
	    }
	}

}

export namespace llms {
	
	export class FunctionCall {
	    name: string;
	    arguments: string;
	
	    static createFrom(source: any = {}) {
	        return new FunctionCall(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.arguments = source["arguments"];
	    }
	}

}

export namespace model {
	
	export class HttpResult {
	    code: number;
	    msg: string;
	    data: any;
	
	    static createFrom(source: any = {}) {
	        return new HttpResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	        this.data = source["data"];
	    }
	}
	export class OrderItem {
	    column: string;
	    asc: boolean;
	
	    static createFrom(source: any = {}) {
	        return new OrderItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.column = source["column"];
	        this.asc = source["asc"];
	    }
	}
	export class PageInfo {
	    current: number;
	    size: number;
	    orders: OrderItem[];
	
	    static createFrom(source: any = {}) {
	        return new PageInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.current = source["current"];
	        this.size = source["size"];
	        this.orders = this.convertValues(source["orders"], OrderItem);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace sql {
	
	export class DB {
	
	
	    static createFrom(source: any = {}) {
	        return new DB(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}

}

