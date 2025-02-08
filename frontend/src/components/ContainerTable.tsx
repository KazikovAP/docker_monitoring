import React, { useEffect, useState } from 'react';
import { Table, Spin, Alert } from 'antd';
import { fetchContainerData } from '../services/api';

interface ContainerData {
  id: string;
  ip_address: string;
  ping_time: number;
  last_success: string;
}

const ContainerTable: React.FC = () => {
  const [data, setData] = useState<ContainerData[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const getData = async () => {
      try {
        const result = await fetchContainerData();
        setData(result);
      } catch (err) {
        setError('Не удалось загрузить данные.');
      } finally {
        setLoading(false);
      }
    };

    getData();
    const interval = setInterval(getData, 5000);

    return () => clearInterval(interval);
  }, []);

  if (loading) return <Spin tip="Загрузка..." />;
  if (error) return <Alert message={error} type="error" />;

  const columns = [
    {
      title: 'IP адрес',
      dataIndex: 'ip_address',
      key: 'ip_address',
    },
    {
      title: 'Время пинга (мс)',
      dataIndex: 'ping_time',
      key: 'ping_time',
    },
    {
      title: 'Последняя успешная попытка',
      dataIndex: 'last_success',
      key: 'last_success',
    },
  ];

  return <Table dataSource={data} columns={columns} rowKey="id" />;
};

export default ContainerTable;
