export namespace app {
	
	export class LCUStatus {
	    connected: boolean;
	    port?: string;
	    authToken?: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new LCUStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.connected = source["connected"];
	        this.port = source["port"];
	        this.authToken = source["authToken"];
	        this.error = source["error"];
	    }
	}

}

export namespace config {
	
	export class Config {
	    riot_api_key: string;
	    region: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.riot_api_key = source["riot_api_key"];
	        this.region = source["region"];
	    }
	}

}

export namespace lcu {
	
	export class RerollPoints {
	    currentPoints: number;
	    maxRolls: number;
	    numberOfRolls: number;
	    pointsCostToRoll: number;
	    pointsToReroll: number;
	
	    static createFrom(source: any = {}) {
	        return new RerollPoints(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.currentPoints = source["currentPoints"];
	        this.maxRolls = source["maxRolls"];
	        this.numberOfRolls = source["numberOfRolls"];
	        this.pointsCostToRoll = source["pointsCostToRoll"];
	        this.pointsToReroll = source["pointsToReroll"];
	    }
	}
	export class CurrentSummoner {
	    accountId: number;
	    displayName: string;
	    gameName: string;
	    tagLine: string;
	    internalName: string;
	    nameChangeFlag: boolean;
	    percentCompleteForNextLevel: number;
	    profileIconId: number;
	    puuid: string;
	    rerollPoints: RerollPoints;
	    summonerId: number;
	    summonerLevel: number;
	    xpSinceLastLevel: number;
	    xpUntilNextLevel: number;
	
	    static createFrom(source: any = {}) {
	        return new CurrentSummoner(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.accountId = source["accountId"];
	        this.displayName = source["displayName"];
	        this.gameName = source["gameName"];
	        this.tagLine = source["tagLine"];
	        this.internalName = source["internalName"];
	        this.nameChangeFlag = source["nameChangeFlag"];
	        this.percentCompleteForNextLevel = source["percentCompleteForNextLevel"];
	        this.profileIconId = source["profileIconId"];
	        this.puuid = source["puuid"];
	        this.rerollPoints = this.convertValues(source["rerollPoints"], RerollPoints);
	        this.summonerId = source["summonerId"];
	        this.summonerLevel = source["summonerLevel"];
	        this.xpSinceLastLevel = source["xpSinceLastLevel"];
	        this.xpUntilNextLevel = source["xpUntilNextLevel"];
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

export namespace lol {
	
	export class ChampionMasteryInfo {
	    championId: number;
	    championLevel: number;
	    championPoints: number;
	    championPointsSinceLastLevel: number;
	    championPointsUntilNextLevel: number;
	    chestGranted: boolean;
	    lastPlayTime: number;
	    tokensEarned: number;
	    summonerId: string;
	
	    static createFrom(source: any = {}) {
	        return new ChampionMasteryInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.championId = source["championId"];
	        this.championLevel = source["championLevel"];
	        this.championPoints = source["championPoints"];
	        this.championPointsSinceLastLevel = source["championPointsSinceLastLevel"];
	        this.championPointsUntilNextLevel = source["championPointsUntilNextLevel"];
	        this.chestGranted = source["chestGranted"];
	        this.lastPlayTime = source["lastPlayTime"];
	        this.tokensEarned = source["tokensEarned"];
	        this.summonerId = source["summonerId"];
	    }
	}
	export class RankedInfo {
	    queueType: string;
	    tier: string;
	    rank: string;
	    leaguePoints: number;
	    wins: number;
	    losses: number;
	    hotStreak: boolean;
	    veteran: boolean;
	    freshBlood: boolean;
	    inactive: boolean;
	    summonerId: string;
	    summonerName: string;
	    puuid: string;
	
	    static createFrom(source: any = {}) {
	        return new RankedInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.queueType = source["queueType"];
	        this.tier = source["tier"];
	        this.rank = source["rank"];
	        this.leaguePoints = source["leaguePoints"];
	        this.wins = source["wins"];
	        this.losses = source["losses"];
	        this.hotStreak = source["hotStreak"];
	        this.veteran = source["veteran"];
	        this.freshBlood = source["freshBlood"];
	        this.inactive = source["inactive"];
	        this.summonerId = source["summonerId"];
	        this.summonerName = source["summonerName"];
	        this.puuid = source["puuid"];
	    }
	}
	export class LeagueListInfo {
	    tier: string;
	    leagueId: string;
	    queue: string;
	    name: string;
	    entries: RankedInfo[];
	
	    static createFrom(source: any = {}) {
	        return new LeagueListInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tier = source["tier"];
	        this.leagueId = source["leagueId"];
	        this.queue = source["queue"];
	        this.name = source["name"];
	        this.entries = this.convertValues(source["entries"], RankedInfo);
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
	
	export class SummonerInfo {
	    id: string;
	    accountId: string;
	    puuid: string;
	    gameName: string;
	    tagLine: string;
	    profileIconId: number;
	    summonerLevel: number;
	    revisionDate: number;
	
	    static createFrom(source: any = {}) {
	        return new SummonerInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.accountId = source["accountId"];
	        this.puuid = source["puuid"];
	        this.gameName = source["gameName"];
	        this.tagLine = source["tagLine"];
	        this.profileIconId = source["profileIconId"];
	        this.summonerLevel = source["summonerLevel"];
	        this.revisionDate = source["revisionDate"];
	    }
	}

}

