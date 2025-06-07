import { useState } from 'react';

interface SalesData {
  date: string;
  orders: number;
  purchases: number;
}

interface SalesChartProps {
  data: SalesData[];
  isLoading?: boolean;
}

export default function SalesChart({ data, isLoading = false }: SalesChartProps) {
  const [activeTab, setActiveTab] = useState<'orders' | 'purchases'>('orders');
  
  if (isLoading) {
    return (
      <div className="bg-white rounded-lg p-4 shadow-sm border border-gray-100">
        <div className="animate-pulse">
          <div className="flex space-x-2 mb-4">
            <div className="h-8 bg-gray-200 rounded w-20"></div>
            <div className="h-8 bg-gray-200 rounded w-20"></div>
          </div>
          <div className="h-48 bg-gray-200 rounded"></div>
        </div>
      </div>
    );
  }

  const maxValue = Math.max(...data.map(d => Math.max(d.orders, d.purchases)));
  const chartData = data.map(d => activeTab === 'orders' ? d.orders : d.purchases);
  
  return (
    <div className="bg-white rounded-lg p-4 shadow-sm border border-gray-100">
      {/* Tabs */}
      <div className="flex space-x-1 mb-4 bg-gray-100 rounded-lg p-1">
        <button
          onClick={() => setActiveTab('orders')}
          className={`flex-1 py-2 px-3 text-sm font-medium rounded-md transition-colors ${
            activeTab === 'orders'
              ? 'bg-white text-gray-900 shadow-sm'
              : 'text-gray-500 hover:text-gray-700'
          }`}
        >
          Заказали
        </button>
        <button
          onClick={() => setActiveTab('purchases')}
          className={`flex-1 py-2 px-3 text-sm font-medium rounded-md transition-colors ${
            activeTab === 'purchases'
              ? 'bg-white text-gray-900 shadow-sm'
              : 'text-gray-500 hover:text-gray-700'
          }`}
        >
          Выкупили
        </button>
      </div>

      {/* Chart */}
      <div className="h-48 relative">
        <svg className="w-full h-full" viewBox="0 0 400 200" preserveAspectRatio="none">
          {/* Grid lines */}
          {[0, 1, 2, 3, 4].map(i => (
            <line
              key={i}
              x1="0"
              y1={i * 40}
              x2="400"
              y2={i * 40}
              stroke="#f3f4f6"
              strokeWidth="1"
            />
          ))}
          
          {/* Chart line */}
          <polyline
            fill="none"
            stroke="#374151"
            strokeWidth="2"
            points={chartData
              .map((value, index) => {
                const x = (index / (chartData.length - 1)) * 400;
                const y = 200 - (value / maxValue) * 180;
                return `${x},${y}`;
              })
              .join(' ')}
          />
          
          {/* Data points */}
          {chartData.map((value, index) => {
            const x = (index / (chartData.length - 1)) * 400;
            const y = 200 - (value / maxValue) * 180;
            return (
              <circle
                key={index}
                cx={x}
                cy={y}
                r="3"
                fill="#374151"
              />
            );
          })}
        </svg>
        
        {/* X-axis labels */}
        <div className="absolute bottom-0 left-0 right-0 flex justify-between text-xs text-gray-500 mt-2">
          {data.map((item, index) => (
            <span key={index}>{new Date(item.date).toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit' })}</span>
          ))}
        </div>
      </div>
    </div>
  );
} 