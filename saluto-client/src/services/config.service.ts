import { Injectable } from "@angular/core";

@Injectable({
    providedIn: 'root'
})
export class ConfigService {
    APIRoot = "https://saluto.site/api/";
    // APIRoot = "http://localhost:8080/api/";
}