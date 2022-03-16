import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../environments/environment';
import { Realtime } from 'ably';

@Injectable({
  providedIn: 'root'
})
export class WalletService {
  client: Realtime = new Realtime(`${environment.API_KEY_COIN}`);

  constructor(private httpClient: HttpClient) {

   }

  getBitcoinPrice() {
    let chanName = '[product:ably-coindesk/bitcoin]bitcoin:usd';
    let channel = this.client.channels.get(chanName);
    return channel;
  }


  createPlayer(data: string) {
    return this.httpClient.post(`${environment.gateway}/start`, {Name: data});
  }

  getWalletBalance(player: string) {
    return this.httpClient.get(`${environment.gateway}/balance/${player}`);
  }

  buyBitcoin(price: number, quantity: number, name: string) {
    let data = {
      Player: name,
      Quantity: quantity,
      CurrentCoinPrice: +price
    }
    return this.httpClient.post(`${environment.gateway}/buy`, data);
  }

  sellBitcoin(price: number, quantity: number, name: string) {
    let data = {
      Player: name,
      Quantity: quantity,
      CurrentCoinPrice: +price
    }
    return this.httpClient.post(`${environment.gateway}/sell`, data);
  }

}
