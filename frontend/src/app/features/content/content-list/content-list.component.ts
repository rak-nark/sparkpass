import { Component, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ContentService } from '../content.service';
import { Router } from '@angular/router';
import { ContentCardComponent } from '../content-card/content-card.component';
import { AuthService } from '../../auth/auth.service'; // Suponiendo que tienes un servicio de autenticación

@Component({
  selector: 'app-content-list',
  standalone: true,
  imports: [CommonModule, ContentCardComponent],
  templateUrl: './content-list.component.html',
  styleUrls: ['./content-list.component.scss']
})

export class ContentListComponent implements OnInit {
  private contentService = inject(ContentService);
  private authService = inject(AuthService); // Servicio de autenticación
  private router = inject(Router);

  contentItems: any[] = [];
  isLoading = true;
  errorMessage = '';
  isCreator: boolean = false;

  ngOnInit(): void {
    this.checkIfCreator();
    this.loadContent();
  }

  checkIfCreator(): void {
    // Aquí verificamos si el usuario es creador (basado en la lógica de tu backend)
    const user = this.authService.getUser(); // Suponiendo que tienes un método que retorna el usuario logueado
    this.isCreator = user && user.isCreator;
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
