package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

// MailgunService handles email sending via Mailgun
type MailgunService struct {
	Client *mailgun.MailgunImpl
}

// Mailer interface for email operations
type Mailer interface {
	SendBookingConfirmation(bookingEmail, organizerName, bookingID, eventType, bookingDate string, guestCount int, totalPrice float64) (string, error)
	SendBookingStatusUpdate(bookingEmail, organizerName, bookingID, eventType, newStatus string) (string, error)
	SendBookingCancellation(bookingEmail, organizerName, bookingID, eventType, bookingDate string) (string, error)
}

// NewMailgunService creates a new Mailgun service instance
func NewMailgunService() *MailgunService {
	domain := os.Getenv("MG_DOMAIN")
	apiKey := os.Getenv("MG_API_KEY")

	if domain == "" || apiKey == "" {
		log.Printf("Mailgun credentials not found in environment variables")
		return nil
	}

	return &MailgunService{
		Client: mailgun.NewMailgun(domain, apiKey),
	}
}

// SendBookingConfirmation sends email when a new booking is created
func (m *MailgunService) SendBookingConfirmation(bookingEmail, organizerName, bookingID, eventType, bookingDate string, guestCount int, totalPrice float64) (string, error) {
	if m.Client == nil {
		return "", fmt.Errorf("mail service not initialized")
	}

	emailFrom := os.Getenv("MG_EMAIL_FROM")
	if emailFrom == "" {
		emailFrom = "noreply@" + os.Getenv("MG_DOMAIN")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Create message
	subject := fmt.Sprintf("Booking Confirmation - %s (%s)", eventType, bookingID)
	htmlBody := fmt.Sprintf(`
		<h2>Booking Confirmation</h2>
		<p>Dear %s,</p>
		<p>Your booking has been confirmed! Here are your booking details:</p>
		
		<table style="border-collapse: collapse; width: 100%%;">
			<tr>
				<td style="padding: 8px; border: 1px solid #ddd;"><strong>Booking ID:</strong></td>
				<td style="padding: 8px; border: 1px solid #ddd;">%s</td>
			</tr>
			<tr>
				<td style="padding: 8px; border: 1px solid #ddd;"><strong>Event Type:</strong></td>
				<td style="padding: 8px; border: 1px solid #ddd;">%s</td>
			</tr>
			<tr>
				<td style="padding: 8px; border: 1px solid #ddd;"><strong>Date:</strong></td>
				<td style="padding: 8px; border: 1px solid #ddd;">%s</td>
			</tr>
			<tr>
				<td style="padding: 8px; border: 1px solid #ddd;"><strong>Guest Count:</strong></td>
				<td style="padding: 8px; border: 1px solid #ddd;">%d</td>
			</tr>
			<tr>
				<td style="padding: 8px; border: 1px solid #ddd;"><strong>Total Price:</strong></td>
				<td style="padding: 8px; border: 1px solid #ddd;">$%.2f</td>
			</tr>
		</table>
		
		<p>Please contact us if you have any questions or need to make changes to your booking.</p>
		<p>Thank you for choosing our venue!</p>
		
		<br>
		<p>Best regards,<br>Venue Management Team</p>
	`, organizerName, bookingID, eventType, bookingDate, guestCount, totalPrice)

	// Create message
	message := m.Client.NewMessage(
		emailFrom,
		subject,
		"", // Empty text body since we're sending HTML
		bookingEmail,
	)
	message.SetHtml(htmlBody)

	// Send email
	messageID, _, err := m.Client.Send(ctx, message)
	if err != nil {
		log.Printf("Failed to send booking confirmation email: %v", err)
		return "", err
	}

	log.Printf("Booking confirmation email sent to %s, Message ID: %s", bookingEmail, messageID)
	return messageID, nil
}

// SendBookingStatusUpdate sends email when booking status changes
func (m *MailgunService) SendBookingStatusUpdate(bookingEmail, organizerName, bookingID, eventType, newStatus string) (string, error) {
	if m.Client == nil {
		return "", fmt.Errorf("mail service not initialized")
	}

	emailFrom := os.Getenv("MG_EMAIL_FROM")
	if emailFrom == "" {
		emailFrom = "noreply@" + os.Getenv("MG_DOMAIN")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Create message
	subject := fmt.Sprintf("Booking Status Update - %s (%s)", eventType, bookingID)
	htmlBody := fmt.Sprintf(`
		<h2>Booking Status Update</h2>
		<p>Dear %s,</p>
		<p>Your booking status has been updated:</p>
		
		<table style="border-collapse: collapse; width: 100%%;">
			<tr>
				<td style="padding: 8px; border: 1px solid #ddd;"><strong>Booking ID:</strong></td>
				<td style="padding: 8px; border: 1px solid #ddd;">%s</td>
			</tr>
			<tr>
				<td style="padding: 8px; border: 1px solid #ddd;"><strong>Event Type:</strong></td>
				<td style="padding: 8px; border: 1px solid #ddd;">%s</td>
			</tr>
			<tr>
				<td style="padding: 8px; border: 1px solid #ddd;"><strong>New Status:</strong></td>
				<td style="padding: 8px; border: 1px solid #ddd;">%s</td>
			</tr>
		</table>
		
		<p>Status Description:</p>
		<ul>
			<li><strong>Pending:</strong> Your booking is awaiting confirmation</li>
			<li><strong>Confirmed:</strong> Your booking has been confirmed and is scheduled</li>
			<li><strong>Completed:</strong> Your event has been successfully conducted</li>
			<li><strong>Cancelled:</strong> Your booking has been cancelled</li>
		</ul>
		
		<p>If you have any questions about this status change, please contact us immediately.</p>
		
		<br>
		<p>Best regards,<br>Venue Management Team</p>
	`, organizerName, bookingID, eventType, newStatus)

	// Create message
	message := m.Client.NewMessage(
		emailFrom,
		subject,
		"", // Empty text body since we're sending HTML
		bookingEmail,
	)
	message.SetHtml(htmlBody)

	// Send email
	messageID, _, err := m.Client.Send(ctx, message)
	if err != nil {
		log.Printf("Failed to send booking status update email: %v", err)
		return "", err
	}

	log.Printf("Booking status update email sent to %s, Message ID: %s", bookingEmail, messageID)
	return messageID, nil
}

// SendBookingCancellation sends email when booking is cancelled
func (m *MailgunService) SendBookingCancellation(bookingEmail, organizerName, bookingID, eventType, bookingDate string) (string, error) {
	if m.Client == nil {
		return "", fmt.Errorf("mail service not initialized")
	}

	emailFrom := os.Getenv("MG_EMAIL_FROM")
	if emailFrom == "" {
		emailFrom = "noreply@" + os.Getenv("MG_DOMAIN")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Create message
	subject := fmt.Sprintf("Booking Cancelled - %s (%s)", eventType, bookingID)
	htmlBody := fmt.Sprintf(`
		<h2>Booking Cancellation</h2>
		<p>Dear %s,</p>
		<p>Your booking has been cancelled. Here are the details:</p>
		
		<table style="border-collapse: collapse; width: 100%%;">
			<tr>
				<td style="padding: 8px; border: 1px solid #ddd;"><strong>Booking ID:</strong></td>
				<td style="padding: 8px; border: 1px solid #ddd;">%s</td>
			</tr>
			<tr>
				<td style="padding: 8px; border: 1px solid #ddd;"><strong>Event Type:</strong></td>
				<td style="padding: 8px; border: 1px solid #ddd;">%s</td>
			</tr>
			<tr>
				<td style="padding: 8px; border: 1px solid #ddd;"><strong>Original Date:</strong></td>
				<td style="padding: 8px; border: 1px solid #ddd;">%s</td>
			</tr>
		</table>
		
		<p>We're sorry to see your booking cancelled. If this was done in error or if you'd like to reschedule, please contact us as soon as possible.</p>
		
		<p>If you have already made any payments, please contact us regarding the refund process.</p>
		
		<br>
		<p>Best regards,<br>Venue Management Team</p>
	`, organizerName, bookingID, eventType, bookingDate)

	// Create message
	message := m.Client.NewMessage(
		emailFrom,
		subject,
		"", // Empty text body since we're sending HTML
		bookingEmail,
	)
	message.SetHtml(htmlBody)

	// Send email
	messageID, _, err := m.Client.Send(ctx, message)
	if err != nil {
		log.Printf("Failed to send booking cancellation email: %v", err)
		return "", err
	}

	log.Printf("Booking cancellation email sent to %s, Message ID: %s", bookingEmail, messageID)
	return messageID, nil
}
