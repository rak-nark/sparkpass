import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { TruncatePipe } from '../../../shared/pipes/truncate.pipe';

@Component({
  selector: 'app-content-card',
  standalone: true,
  imports: [CommonModule, RouterLink, TruncatePipe],
  template: `
    <div class="content-card">
      <h3>{{ content.title }}</h3>
      <img [src]="content.thumbnail" [alt]="content.title">
      <p>{{ content.description | truncate:100 }}</p>
      <a [routerLink]="['/content', content.id]" class="view-button">Ver más</a>
    </div>
  `,
  styleUrls: ['./content-card.component.scss']
})
export class ContentCardComponent {
  @Input() content: any;
}