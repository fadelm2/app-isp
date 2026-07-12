import axios from 'axios';

const api = axios.create({
    baseURL: '/api', // Proxied by Vite to localhost:8080
    headers: {
        'Content-Type': 'application/json',
    },
});

// Auto attach authorization token if present in localStorage/cookie
api.interceptors.request.use((config) => {
    const token = localStorage.getItem('token') || document.cookie.replace(/(?:(?:^|.*;\s*)token\s*\=\s*([^;]*).*$)|^.*$/, "$1");
    if (token) {
        config.headers['Authorization'] = token;
    }
    return config;
}, (error) => {
    return Promise.reject(error);
});

export const authService = {
    login: (data) => api.post('/users/_Login', data),
    logout: () => api.delete('/users'),
    getCurrentUser: () => api.get('/users/_current'),
    updateProfile: (data) => api.patch('/users/_current', data),
};

export const publicService = {
    getPackages: () => api.get('/packages'),
    register: (formData) => api.post('/registrations', formData, {
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    })
};

export const adminService = {
    checkAccess: () => api.get('/admin/'),
    
    // Dashboard Stats
    getDashboardStats: () => api.get('/admin/dashboard'),

    // Internet Packages
    getPackages: () => api.get('/admin/packages'),
    createPackage: (data) => api.post('/admin/packages', data),
    updatePackage: (id, data) => api.patch(`/admin/packages/${id}`, data),
    deletePackage: (id) => api.delete(`/admin/packages/${id}`),

    // Customer Registrations
    getRegistrations: () => api.get('/admin/registrations'),
    updateRegistrationStatus: (id, status) => api.patch(`/admin/registrations/${id}/status`, { status }),

    // Customers
    getCustomers: () => api.get('/admin/customers'),
    getCustomer: (id) => api.get(`/admin/customers/${id}`),
    updateCustomer: (id, data) => api.patch(`/admin/customers/${id}`, data),
    suspendCustomer: (id, notes) => api.post(`/admin/customers/${id}/_suspend`, { notes }),
    unsuspendCustomer: (id, notes) => api.post(`/admin/customers/${id}/_unsuspend`, { notes }),
    terminateCustomer: (id, notes) => api.post(`/admin/customers/${id}/_terminate`, { notes }),
    getCustomerHistory: (id) => api.get(`/admin/customers/${id}/history`),

    // Routers
    getRouters: () => api.get('/admin/routers'),
    createRouter: (data) => api.post('/admin/routers', data),
    deleteRouter: (id) => api.delete(`/admin/routers/${id}`),

    // Invoices
    getInvoices: () => api.get('/admin/invoices'),
    createInvoice: (data) => api.post('/admin/invoices', data),
    getInvoice: (id) => api.get(`/admin/invoices/${id}`),
};

export const customerService = {
    getProfile: () => api.get('/customer/me'),
    getSnapToken: (invoiceId) => api.get(`/customer/invoices/${invoiceId}/pay`),
};

export default api;
