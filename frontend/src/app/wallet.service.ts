import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../environments/environment';
import { Realtime } from 'ably';
import { Observable, Observer } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class WalletService {
  observer!: Observer<number>;
  ably = new Realtime('VbpYdQ.79UpAA:-lnejxoRLhS_hDPgNrE5XqweLrsLdH0vMZwSQtaKlLI');
  chanName = '[product:ably-coindesk/bitcoin]bitcoin:usd';
  channel$ = this.ably.channels.get(this.chanName);
  // public messages$ = this.channel$.

  constructor(private httpClient: HttpClient) { }

  getCurrentPrice() {
    // let ably = new Realtime('VbpYdQ.79UpAA:-lnejxoRLhS_hDPgNrE5XqweLrsLdH0vMZwSQtaKlLI');
    // let chanName = '[product:ably-coindesk/bitcoin]bitcoin:usd';
    // let channel = ably.channels.get(chanName);




    // channel.subscribe((message: any): void => {
    //   console.log('Received msg', message.data);
    //   return this.observer.next(message.data);
    // });
    // return this.createObservable();
  }

  // createObservable() {
  //   return new Observable((observer) => (this.observer = observer));
  // }

  getWalletBalance() {
    return this.httpClient.get(environment.gateway + '/play/');
  }
}
