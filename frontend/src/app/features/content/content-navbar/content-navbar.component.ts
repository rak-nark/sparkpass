// content-navbar.component.ts
import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../../auth/auth.service';
import { UserService } from '../../../services/user.service';
import { NotificationService } from '../../../services/notification.service'; // Añade esta importación
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { RouterLink, RouterLinkActive } from '@angular/router';

@Component({
  selector: 'app-content-navbar',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    RouterLink
],
  templateUrl: './content-navbar.component.html',
  styleUrls: ['./content-navbar.component.scss']
})
export class ContentNavbarComponent implements OnInit {
  searchQuery: string = '';
  showProfileMenu: boolean = false;
  showNotifications: boolean = false; // Añadido
  userProfile: any;
  isCreator: boolean = false;
  unreadNotifications: number = 0; // Añadido
  notifications: any[] = []; // Añadido

  constructor(
    private router: Router,
    public authService: AuthService,
    private userService: UserService,
    private notificationService: NotificationService // Añadido al constructor
  ) {}

  ngOnInit(): void {
    this.loadUserProfile();
    this.loadNotifications();
  }

  loadUserProfile(): void {
    this.userService.getCurrentUser().subscribe({
      next: (user) => {
        this.userProfile = user;
        this.isCreator = user?.isCreator || false;
      },
      error: (err) => console.error('Error loading user profile', err)
    });
  }

  loadNotifications(): void {
    this.notificationService.getNotifications().subscribe({
      next: (notifications) => {
        this.notifications = notifications;
        this.unreadNotifications = notifications.filter(n => !n.read).length;
      },
      error: (err) => console.error('Error loading notifications', err)
    });
  }

  onSearch(): void {
    if (this.searchQuery.trim()) {
      this.router.navigate(['/content/search'], { 
        queryParams: { q: this.searchQuery } 
      });
    }
  }

  toggleNotifications(): void {
    this.showNotifications = !this.showNotifications;
    if (this.showNotifications) {
      this.showProfileMenu = false;
      if (this.unreadNotifications > 0) {
        this.notificationService.markAsRead().subscribe(() => {
          this.unreadNotifications = 0;
        });
      }
    }
  }

  toggleProfileMenu(): void {
    this.showProfileMenu = !this.showProfileMenu;
    if (this.showProfileMenu) {
      this.showNotifications = false;
    }
  }

  logout(): void {
    this.authService.logout();
    this.router.navigate(['/login']);
  }
}