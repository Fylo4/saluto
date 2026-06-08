import { Injectable } from "@angular/core";

@Injectable({
    providedIn: 'root'
})
export class ConfigService {
    APIRoot = "http://34.196.129.60:8080/api/";
    // APIRoot = "http://localhost:8080/api/";
}