import { HttpClient } from "@angular/common/http";
import { inject, Injectable } from "@angular/core";

const API_ROOT = "http://34.196.129.60:8080/api/";

@Injectable({
    providedIn: 'root'  
})
export class ApiService {
    private http = inject(HttpClient);

    public testServer() {
        return this.http.get<{message: string}>(API_ROOT)
    }
}