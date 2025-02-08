import React from 'react';
import { Layout, Typography } from 'antd';
import ContainerTable from './components/ContainerTable';

const { Header, Content } = Layout;
const { Title } = Typography;

const App: React.FC = () => {
  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ background: '#001529' }}>
        <Title style={{ color: '#fff', textAlign: 'center' }}>Мониторинг Контейнеров</Title>
      </Header>
      <Content style={{ padding: '20px' }}>
        <ContainerTable />
      </Content>
    </Layout>
  );
};

export default App;
