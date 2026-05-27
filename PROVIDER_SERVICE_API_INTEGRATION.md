# Provider Service CRUD API - Frontend Integration Guide

## Overview
New provider service management API endpoints have been created to allow providers to manage their services (create, read, update, delete, toggle availability).

## API Endpoints

### 1. Get Provider Services
**Endpoint:** `GET /api/v1/provider/services`

**Authentication:** Required (Bearer token)

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "provider_id": "uuid",
      "title": "Hotel Room Booking",
      "description": "Luxury room with ocean view",
      "category_id": "category-uuid",
      "price_type": "fixed",
      "base_price": 50000,
      "duration": 60,
      "features": ["WiFi", "AC", "Breakfast"],
      "images": ["https://example.com/image1.jpg"],
      "is_active": true,
      "is_available": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "message": "Provider services retrieved"
}
```

---

### 2. Create Provider Service
**Endpoint:** `POST /api/v1/provider/services`

**Authentication:** Required (Bearer token)

**Request Body:**
```json
{
  "title": "Hotel Room Booking",
  "description": "Luxury room with ocean view",
  "category_id": "category-uuid",
  "price_type": "fixed",
  "base_price": 50000,
  "duration": 60,
  "features": ["WiFi", "AC", "Breakfast"],
  "images": ["https://example.com/image1.jpg"]
}
```

**Required Fields:**
- `title` (string)
- `price_type` (string) - "fixed", "hourly", "per_item", "custom"
- `base_price` (number)

**Optional Fields:**
- `description` (string)
- `category_id` (string)
- `duration` (number) - in minutes
- `features` (array of strings)
- `images` (array of strings)

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "provider_id": "uuid",
    "title": "Hotel Room Booking",
    "description": "Luxury room with ocean view",
    "category_id": "category-uuid",
    "price_type": "fixed",
    "base_price": 50000,
    "duration": 60,
    "features": ["WiFi", "AC", "Breakfast"],
    "images": ["https://example.com/image1.jpg"],
    "is_active": true,
    "is_available": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "message": "Service created successfully"
}
```

---

### 3. Update Provider Service
**Endpoint:** `PUT /api/v1/provider/services/:id`

**Authentication:** Required (Bearer token)

**URL Parameters:**
- `id` (string) - Service UUID

**Request Body:**
```json
{
  "title": "Updated Service Title",
  "description": "Updated description",
  "category_id": "new-category-uuid",
  "price_type": "hourly",
  "base_price": 75000,
  "duration": 90,
  "features": ["WiFi", "AC", "Breakfast", "Pool"],
  "images": ["https://example.com/image1.jpg", "https://example.com/image2.jpg"],
  "is_active": true,
  "is_available": false
}
```

**Note:** All fields are optional. Only provided fields will be updated.

**Response:**
```json
{
  "success": true,
  "message": "Service updated successfully"
}
```

---

### 4. Delete Provider Service
**Endpoint:** `DELETE /api/v1/provider/services/:id`

**Authentication:** Required (Bearer token)

**URL Parameters:**
- `id` (string) - Service UUID

**Response:**
```json
{
  "success": true,
  "message": "Service deleted successfully"
}
```

---

### 5. Toggle Service Availability
**Endpoint:** `PATCH /api/v1/provider/services/:id/availability`

**Authentication:** Required (Bearer token)

**URL Parameters:**
- `id` (string) - Service UUID

**Request Body:**
```json
{
  "is_available": true
}
```

**Response:**
```json
{
  "success": true,
  "message": "Service availability updated"
}
```

---

## Frontend Integration

### 1. Add to API Service Module

**File:** `lib/api/provider.ts`

```typescript
// Add these methods to providerAPI
export const providerAPI = {
  // ... existing methods ...

  // Service Management
  getServices: () => apiClient.get('/provider/services'),
  createService: (data) => apiClient.post('/provider/services', data),
  updateService: (id, data) => apiClient.put(`/provider/services/${id}`, data),
  deleteService: (id) => apiClient.delete(`/provider/services/${id}`),
  toggleServiceAvailability: (id, isAvailable) => 
    apiClient.patch(`/provider/services/${id}/availability`, { is_available: isAvailable }),
};
```

---

### 2. TypeScript Types

**File:** `types/provider.ts`

```typescript
export interface ProviderService {
  id: string;
  provider_id: string;
  title: string;
  description: string;
  category_id: string;
  price_type: 'fixed' | 'hourly' | 'per_item' | 'custom';
  base_price: number;
  duration: number;
  features: string[];
  images: string[];
  is_active: boolean;
  is_available: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateServiceRequest {
  title: string;
  description?: string;
  category_id?: string;
  price_type: 'fixed' | 'hourly' | 'per_item' | 'custom';
  base_price: number;
  duration?: number;
  features?: string[];
  images?: string[];
}

export interface UpdateServiceRequest {
  title?: string;
  description?: string;
  category_id?: string;
  price_type?: 'fixed' | 'hourly' | 'per_item' | 'custom';
  base_price?: number;
  duration?: number;
  features?: string[];
  images?: string[];
  is_active?: boolean;
  is_available?: boolean;
}
```

---

### 3. React Hook Example

**File:** `hooks/useProviderServices.ts`

```typescript
import { useState, useEffect } from 'react';
import { providerAPI } from '@/lib/api/provider';
import type { ProviderService, CreateServiceRequest, UpdateServiceRequest } from '@/types/provider';

export function useProviderServices() {
  const [services, setServices] = useState<ProviderService[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchServices = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await providerAPI.getServices();
      setServices(response.data.data);
    } catch (err) {
      setError('Failed to fetch services');
    } finally {
      setLoading(false);
    }
  };

  const createService = async (data: CreateServiceRequest) => {
    setLoading(true);
    setError(null);
    try {
      const response = await providerAPI.createService(data);
      setServices([...services, response.data.data]);
      return response.data.data;
    } catch (err) {
      setError('Failed to create service');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const updateService = async (id: string, data: UpdateServiceRequest) => {
    setLoading(true);
    setError(null);
    try {
      await providerAPI.updateService(id, data);
      await fetchServices(); // Refresh list
    } catch (err) {
      setError('Failed to update service');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const deleteService = async (id: string) => {
    setLoading(true);
    setError(null);
    try {
      await providerAPI.deleteService(id);
      setServices(services.filter(s => s.id !== id));
    } catch (err) {
      setError('Failed to delete service');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const toggleAvailability = async (id: string, isAvailable: boolean) => {
    setLoading(true);
    setError(null);
    try {
      await providerAPI.toggleServiceAvailability(id, isAvailable);
      await fetchServices(); // Refresh list
    } catch (err) {
      setError('Failed to update availability');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchServices();
  }, []);

  return {
    services,
    loading,
    error,
    fetchServices,
    createService,
    updateService,
    deleteService,
    toggleAvailability,
  };
}
```

---

### 4. Component Integration Example

**File:** `pages/provider/services/index.tsx`

```typescript
import { useProviderServices } from '@/hooks/useProviderServices';

export default function ProviderServicesPage() {
  const {
    services,
    loading,
    error,
    createService,
    updateService,
    deleteService,
    toggleAvailability,
  } = useProviderServices();

  const handleCreateService = async (formData) => {
    try {
      await createService(formData);
      // Show success message
    } catch (err) {
      // Show error message
    }
  };

  const handleToggleAvailability = async (serviceId, currentStatus) => {
    try {
      await toggleAvailability(serviceId, !currentStatus);
    } catch (err) {
      // Show error message
    }
  };

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div>
      <h1>My Services</h1>
      
      {/* Add Service Button */}
      <button onClick={() => setShowCreateModal(true)}>
        Add New Service
      </button>

      {/* Services List */}
      <div className="services-grid">
        {services.map((service) => (
          <div key={service.id} className="service-card">
            <h3>{service.title}</h3>
            <p>{service.description}</p>
            <p>Price: ₦{service.base_price}</p>
            
            {/* Availability Toggle */}
            <button
              onClick={() => handleToggleAvailability(service.id, service.is_available)}
            >
              {service.is_available ? 'Available' : 'Unavailable'}
            </button>

            {/* Edit/Delete Actions */}
            <button onClick={() => handleEdit(service)}>Edit</button>
            <button onClick={() => handleDelete(service.id)}>Delete</button>
          </div>
        ))}
      </div>
    </div>
  );
}
```

---

## Price Types

- **fixed** - One-time payment (e.g., hotel room)
- **hourly** - Per hour rate (e.g., driver service)
- **per_item** - Per unit/item (e.g., food delivery)
- **custom** - Custom pricing structure

---

## Image Upload

For image uploads, you'll need to:

1. Upload images to your file storage service (S3, Cloudinary, etc.)
2. Get the image URLs
3. Pass the URLs array in the `images` field

Example:
```typescript
const handleImageUpload = async (files: File[]) => {
  const uploadedUrls = await Promise.all(
    files.map(file => uploadToStorage(file))
  );
  return uploadedUrls;
};
```

---

## Error Handling

All endpoints return errors in this format:

```json
{
  "success": false,
  "message": "Error message here"
}
```

Common HTTP status codes:
- `401` - Unauthorized (token missing or invalid)
- `400` - Bad request (validation error)
- `404` - Not found (service doesn't exist)
- `500` - Internal server error

---

## Testing

You can test these endpoints using:

```bash
# Get services
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/v1/provider/services

# Create service
curl -X POST \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Test Service","price_type":"fixed","base_price":50000}' \
  http://localhost:8080/api/v1/provider/services
```

---

## Summary

**New Endpoints:**
- ✅ `GET /api/v1/provider/services` - List services
- ✅ `POST /api/v1/provider/services` - Create service
- ✅ `PUT /api/v1/provider/services/:id` - Update service
- ✅ `DELETE /api/v1/provider/services/:id` - Delete service
- ✅ `PATCH /api/v1/provider/services/:id/availability` - Toggle availability

**Frontend Tasks:**
1. Add API methods to `lib/api/provider.ts`
2. Add TypeScript types to `types/provider.ts`
3. Create React hook `useProviderServices`
4. Connect service management page to API
5. Implement image upload functionality
6. Add loading states and error handling
