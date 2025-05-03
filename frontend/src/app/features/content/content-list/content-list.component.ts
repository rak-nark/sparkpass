import { Component, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ContentService } from '../content.service';
import { RouterLink } from '@angular/router';
import { ContentCardComponent } from '../content-card/content-card.component';

@Component({
  selector: 'app-content-list',
  standalone: true,
  imports: [CommonModule, ContentCardComponent],
  templateUrl: './content-list.component.html',
  styleUrls: ['./content-list.component.scss']
})
export class ContentListComponent implements OnInit {
  private contentService = inject(ContentService);
  
  contentItems: any[] = [];
  isLoading = true;
  errorMessage = '';

  ngOnInit(): void {
    this.loadContent();
  }

  loadContent(): void {
    this.contentService.getContent().subscribe({
      next: (data) => {
        this.contentItems = data;
        this.isLoading = false;
      },
      error: (err) => {
        this.errorMessage = 'Error al cargar el contenido';
        this.isLoading = false;
        console.error('Error loading content:', err);
      }
    });
  }
}