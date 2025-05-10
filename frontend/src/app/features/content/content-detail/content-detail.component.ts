// content-detail.component.ts
import { Component, inject, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { ContentService } from '../content.service';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-content-detail',
  standalone: true, 
  imports: [CommonModule],
  templateUrl: './content-detail.component.html',
  styleUrls: ['./content-detail.component.scss']
})

export class ContentDetailComponent implements OnInit {
  private contentService = inject(ContentService);
  
  content: any;
  isLoading = true;
  errorMessage = '';

  constructor(private route: ActivatedRoute) {}

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id'); // o 'slug' si usas slug para la URL
    if (id) {
      this.loadContentDetails(id);
    }
  }

  loadContentDetails(id: string): void {
    this.contentService.getContentById(id).subscribe({
      next: (data) => {
        this.content = data;
        this.isLoading = false;
      },
      error: (err) => {
        this.errorMessage = 'Error al cargar el contenido';
        this.isLoading = false;
        console.error('Error loading content details:', err);
      }
    });
  }
}
