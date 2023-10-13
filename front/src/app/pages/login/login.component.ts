import { Component } from '@angular/core';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss'],
})
export class LoginComponent {
  lang: string = 'pt-br';

  ngOnInit(): void {
    this.lang = localStorage.getItem('lang') || 'pt-br';
  }

  changeLang(event: string) {
    localStorage.setItem('lang', event);
    window.location.reload();
  }
}
