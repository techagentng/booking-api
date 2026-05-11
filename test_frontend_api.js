// Test script to verify frontend API integration
const fetch = require('node-fetch');

async function testAPI() {
  try {
    console.log('🧪 Testing Featured Services API...');
    
    const response = await fetch('http://localhost:8080/api/v1/public/explore/featured');
    const data = await response.json();
    
    console.log('📊 Response Status:', response.status);
    console.log('📦 Response Structure:', JSON.stringify(data, null, 2));
    console.log('📋 Data Type:', typeof data.data);
    console.log('📏 Is Array?:', Array.isArray(data.data));
    console.log('📏 Data Length:', data.data ? data.data.length : 'undefined');
    
    if (data.data && Array.isArray(data.data)) {
      console.log('✅ API returns correct array structure');
      console.log('🔄 Testing .map() function...');
      const testMap = data.data.map(item => item.title);
      console.log('✅ .map() works:', testMap);
    } else {
      console.log('❌ API response structure issue');
      console.log('Expected: {success: true, data: [...]}');
      console.log('Got:', data);
    }
    
  } catch (error) {
    console.error('❌ API Test Failed:', error.message);
  }
}

testAPI();
