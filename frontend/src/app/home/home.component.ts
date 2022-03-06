import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { WalletService } from '../wallet.service';
import { Realtime } from 'ably';


@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {
  client: Realtime = new Realtime({
    key: 'VbpYdQ.79UpAA:-lnejxoRLhS_hDPgNrE5XqweLrsLdH0vMZwSQtaKlLI',
    clientId: "hsyd"
  });

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
  members: any = [];


  constructor(public cd: ChangeDetectorRef, private walletService: WalletService) {

  }

  ngOnInit() {
    this.channel.subscribe((message) => {
      let userData = message.data
      const index = this.subscriptions.findIndex(((obj: { Player: any; }) => {
        return obj.Player == userData.Player;
      }));
      if (index === -1) {
        this.subscriptions.push(userData)
      } else {
        this.subscriptions[index] = userData
      }

      this.subscriptions.sort((a: any, b: any) => (a.WalletBalance < b.WalletBalance)? 1: -1)
    });


    // this.updatePresence();
    this.channel.presence.subscribe((presence) => {
      console.log(presence, presence.data, presence.action)
      if (presence.action == "enter") {
        this.members.push(presence.data)
      } else {
        this.members = this.members.filter((member: { name: any; }) => {
          console.log(member.name, presence.data.name, member.name == presence.data.name)
          return member.name != presence.data.name;
        });
      }
    })

    this.getCurrentPrice();
    let check = localStorage.getItem("player")
    let ok = JSON.parse(`${check}`)
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
  }

  createPlayer(): void {
    this.start = true;
    this.walletService.createPlayer(this.playerName).subscribe(() => {
      this.getWallet(this.playerName)
      localStorage.setItem("player", JSON.stringify({ name: this.playerName, start: this.start }))
    })
    this.channel.presence.enter({ name: this.playerName }, (err) => {
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
    this.channel.presence.leave((err: any) => {
      console.log("we're outside", err)
    });
    // this.updatePresence()
    this.start = false;
    localStorage.removeItem("player")
    this.playerName = ""
  }

  // updatePresence() {
  //   this.channel.presence.subscribe((presence) => {
  //     console.log(presence, presence.data, presence.action)
  //     if (presence.action == "enter") {
  //       console.log("enter")
  //       this.members.push(presence.data)
  //     } else {
  //       this.members = this.members.filter((member: { name: any; }) => {
  //         console.log(member.name, presence.data.name, member.name == presence.data.name)
  //         return member.name != presence.data.name;
  //       });
  //     }
  //     console.log(this.members)
  //   })
  // }

}

