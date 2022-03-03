import { ChangeDetectorRef, Component, OnInit} from '@angular/core';
import { WalletService } from '../wallet.service';
import { Realtime } from 'ably';


@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {
  client: Realtime = new Realtime('VbpYdQ.79UpAA:-lnejxoRLhS_hDPgNrE5XqweLrsLdH0vMZwSQtaKlLI');

  portfolio: any = "";
  bitcoinPrice!: number;
  playerName: any = "";
  player: any = "";
  coinWorth!: number;
  quantity: number = 0;
  start: boolean = false;
  chanName = "stock";
  channel = this.client.channels.get(this.chanName);
  view = false;
  subscriptions: any = [];

  constructor(public cd: ChangeDetectorRef, private walletService: WalletService) {

  }



  ngOnInit() {
    this.getCurrentPrice();
    let check = localStorage.getItem("player")
    let ok = JSON.parse(`${check}`)
    console.log(ok)
    if (ok != null) {
      this.playerName = ok.name
      this.start = ok.start
    }
    if (this.start == true) {
      this.getWallet(this.playerName)
    }
  }

  getCurrentPrice(): void {
    this.walletService.getBitcoinPrice().subscribe((data) => {
      this.bitcoinPrice = data.data;
      this.coinWorth = +(this.portfolio.CoinQuantity * this.bitcoinPrice).toFixed(2);
      this.cd.detectChanges();
    })
  }

  getWallet(name: string): void {
    this.walletService.getWalletBalance(name).subscribe((data) => {
      this.portfolio = data;
    });
    this.channel.subscribe((data) => {
      this.subscriptions.push(data.data)
        console.log(this.subscriptions, "subscribed data")
        this.cd.detectChanges()
        console.log(this.subscriptions.length, "len data")
      });
  }

  createPlayer(): void {
    this.start = true;
    this.walletService.createPlayer(this.playerName).subscribe(() => {
      this.getWallet(this.playerName)
      localStorage.setItem("player", JSON.stringify({name: this.playerName, start: this.start}))
      console.log(this.playerName, "set")
    })
    this.channel.presence.enterClient(this.playerName, "playing", (err) => {
      console.log("we're inside")
    });
    this.channel.presence.history((err, presenceSet) => {
      console.log(presenceSet, "set")
    });
  }

  buyCoin(quantity: string, name: string): void {
    console.log(this.bitcoinPrice, quantity)
    this.walletService.buyBitcoin(this.bitcoinPrice, +quantity, name).subscribe((response) => {
      this.getWallet(name)
    });
  }

  sellCoin(quantity: string, name: string): void {
    console.log(this.bitcoinPrice, quantity)
    this.walletService.sellBitcoin(this.bitcoinPrice, +quantity, name).subscribe((response) => {
      this.getWallet(name)
    });

  }

  reset(): void {
    console.log(this.playerName)
    this.start = false;
    localStorage.removeItem("player")
    this.playerName = ""
  }
}

