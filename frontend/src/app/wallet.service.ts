import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../environments/environment';
import { Realtime } from 'ably';

@Injectable({
  providedIn: 'root'
})
export class WalletService {

  constructor(private httpClient: HttpClient) { }

  getWalletBalance() {
    return this.httpClient.get(environment.gateway + '/play/');
  }

  getBitcoinPrice() {
    const ably = new Realtime('VbpYdQ.79UpAA:-lnejxoRLhS_hDPgNrE5XqweLrsLdH0vMZwSQtaKlLI');
    const chanName = '[product:ably-coindesk/bitcoin]bitcoin:usd';
    const channel = ably.channels.get(chanName);

    return channel;
  }

  buyBitcoin(price: number, quantity: number) {
    return this.httpClient.get(environment.gateway + `/play/buy?price=${price}&qty=${quantity}`);
  }
}
