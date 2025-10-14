export namespace models {
	
	export class GraphicsFiles {
	    goalsDirectory: string;
	    technologiesDirectory: string;
	
	    static createFrom(source: any = {}) {
	        return new GraphicsFiles(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.goalsDirectory = source["goalsDirectory"];
	        this.technologiesDirectory = source["technologiesDirectory"];
	    }
	}
	export class InterfaceFiles {
	    goalsGfx: string;
	    countryTechTreeViewGfx: string;
	    otherGfxFiles: string[];
	
	    static createFrom(source: any = {}) {
	        return new InterfaceFiles(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.goalsGfx = source["goalsGfx"];
	        this.countryTechTreeViewGfx = source["countryTechTreeViewGfx"];
	        this.otherGfxFiles = source["otherGfxFiles"];
	    }
	}
	export class ModInfo {
	    basePath: string;
	    name: string;
	    isValid: boolean;
	    validationErrors: string[];
	    nationalFocusFiles: string[];
	    technologyFiles: string[];
	    interfaceFiles: InterfaceFiles;
	    graphicsFiles: GraphicsFiles;
	
	    static createFrom(source: any = {}) {
	        return new ModInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.basePath = source["basePath"];
	        this.name = source["name"];
	        this.isValid = source["isValid"];
	        this.validationErrors = source["validationErrors"];
	        this.nationalFocusFiles = source["nationalFocusFiles"];
	        this.technologyFiles = source["technologyFiles"];
	        this.interfaceFiles = this.convertValues(source["interfaceFiles"], InterfaceFiles);
	        this.graphicsFiles = this.convertValues(source["graphicsFiles"], GraphicsFiles);
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
	export class ValidationResult {
	    isValid: boolean;
	    errors: string[];
	    warnings: string[];
	    modInfo?: ModInfo;
	
	    static createFrom(source: any = {}) {
	        return new ValidationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.isValid = source["isValid"];
	        this.errors = source["errors"];
	        this.warnings = source["warnings"];
	        this.modInfo = this.convertValues(source["modInfo"], ModInfo);
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

