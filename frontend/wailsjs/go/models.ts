export namespace dbstructs {
	
	export class Index {
	    name: string;
	    columns: string[];
	
	    static createFrom(source: any = {}) {
	        return new Index(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.columns = source["columns"];
	    }
	}
	export class Column {
	    columnName: string;
	    data_type: string;
	    not_null: boolean;
	    unique: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Column(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.columnName = source["columnName"];
	        this.data_type = source["data_type"];
	        this.not_null = source["not_null"];
	        this.unique = source["unique"];
	    }
	}
	export class TableMetadata {
	    tableName: string;
	    columns: Column[];
	    primary_key: string[];
	    indexes: Index[];
	    relationships: RelationshipMetadata[];
	
	    static createFrom(source: any = {}) {
	        return new TableMetadata(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tableName = source["tableName"];
	        this.columns = this.convertValues(source["columns"], Column);
	        this.primary_key = source["primary_key"];
	        this.indexes = this.convertValues(source["indexes"], Index);
	        this.relationships = this.convertValues(source["relationships"], RelationshipMetadata);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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

