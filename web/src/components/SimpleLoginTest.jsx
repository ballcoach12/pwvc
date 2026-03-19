import { Box, Container, Typography } from '@mui/material';

const SimpleLoginTest = () => {
  return (
    <Container maxWidth="sm" sx={{ mt: 8 }}>
      <Box textAlign="center">
        <Typography variant="h3" component="h1" gutterBottom color="primary" fontWeight="bold">
          PairWise Login Test
        </Typography>
        <Typography variant="h6" color="text.secondary" gutterBottom>
          This is a simple test to verify React is working!
        </Typography>
        <Typography variant="body1" color="success.main" sx={{ mt: 4, p: 2, bgcolor: 'success.light', borderRadius: 1 }}>
          ✅ SUCCESS: React is rendering correctly!
        </Typography>
      </Box>
    </Container>
  );
};

export default SimpleLoginTest;