import { useState } from 'react';
import TabBar from '../components/TabBar';

export default function ProfilePage() {
  const [user] = useState({
    name: 'Егор Гельман',
    phone: '+7 909 066 1151',
    username: '@egor_gelman1',
    avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1000&q=80'
  });

  const menuItems = [
    {
      id: 'emoji-status',
      icon: (
        <div className="w-8 h-8 bg-yellow-100 rounded-full flex items-center justify-center">
          <span className="text-lg">😊</span>
        </div>
      ),
      title: 'Изменить статус эмодзи',
      hasArrow: false
    },
    {
      id: 'profile-photo',
      icon: (
        <div className="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
          <svg className="w-5 h-5 text-blue-600" fill="currentColor" viewBox="0 0 20 20">
            <path fillRule="evenodd" d="M4 3a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V5a2 2 0 00-2-2H4zm12 12H4l4-8 3 6 2-4 3 6z" clipRule="evenodd" />
          </svg>
        </div>
      ),
      title: 'Изменить фото профиля',
      hasArrow: false
    }
  ];

  const supportItems = [
    {
      id: 'support',
      icon: (
        <div className="w-8 h-8 bg-orange-100 rounded-full flex items-center justify-center">
          <svg className="w-5 h-5 text-orange-600" fill="currentColor" viewBox="0 0 20 20">
            <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-8-3a1 1 0 00-.867.5 1 1 0 11-1.731-1A3 3 0 0113 8a3.001 3.001 0 01-2 2.83V11a1 1 0 11-2 0v-1a1 1 0 011-1 1 1 0 100-2zm0 8a1 1 0 100-2 1 1 0 000 2z" clipRule="evenodd" />
          </svg>
        </div>
      ),
      title: 'LAMODA SUPPORT',
      subtitle: '8',
      hasArrow: true
    }
  ];

  const accountItems = [
    {
      id: 'add-account',
      icon: (
        <div className="w-8 h-8 bg-gray-100 rounded-full flex items-center justify-center">
          <svg className="w-5 h-5 text-gray-600" fill="currentColor" viewBox="0 0 20 20">
            <path fillRule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clipRule="evenodd" />
          </svg>
        </div>
      ),
      title: 'Добавить аккаунт',
      hasArrow: false
    }
  ];

  const otherItems = [
    {
      id: 'my-profile',
      icon: (
        <div className="w-8 h-8 bg-red-100 rounded-full flex items-center justify-center">
          <svg className="w-5 h-5 text-red-600" fill="currentColor" viewBox="0 0 20 20">
            <path fillRule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clipRule="evenodd" />
          </svg>
        </div>
      ),
      title: 'Мой профиль',
      hasArrow: true
    },
    {
      id: 'wallet',
      icon: (
        <div className="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
          <svg className="w-5 h-5 text-blue-600" fill="currentColor" viewBox="0 0 20 20">
            <path d="M4 4a2 2 0 00-2 2v1h16V6a2 2 0 00-2-2H4z" />
            <path fillRule="evenodd" d="M18 9H2v5a2 2 0 002 2h12a2 2 0 002-2V9zM4 13a1 1 0 011-1h1a1 1 0 110 2H5a1 1 0 01-1-1zm5-1a1 1 0 100 2h1a1 1 0 100-2H9z" clipRule="evenodd" />
          </svg>
        </div>
      ),
      title: 'Кошелек',
      hasArrow: true
    },
    {
      id: 'saved-messages',
      icon: (
        <div className="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
          <svg className="w-5 h-5 text-blue-600" fill="currentColor" viewBox="0 0 20 20">
            <path d="M5 4a2 2 0 012-2h6a2 2 0 012 2v14l-5-2.5L5 18V4z" />
          </svg>
        </div>
      ),
      title: 'Сохраненные сообщения',
      hasArrow: true
    },
    {
      id: 'recent-calls',
      icon: (
        <div className="w-8 h-8 bg-green-100 rounded-full flex items-center justify-center">
          <svg className="w-5 h-5 text-green-600" fill="currentColor" viewBox="0 0 20 20">
            <path d="M2 3a1 1 0 011-1h2.153a1 1 0 01.986.836l.74 4.435a1 1 0 01-.54 1.06l-1.548.773a11.037 11.037 0 006.105 6.105l.774-1.548a1 1 0 011.059-.54l4.435.74a1 1 0 01.836.986V17a1 1 0 01-1 1h-2C7.82 18 2 12.18 2 5V3z" />
          </svg>
        </div>
      ),
      title: 'Недавние звонки',
      hasArrow: true
    },
    {
      id: 'devices',
      icon: (
        <div className="w-8 h-8 bg-orange-100 rounded-full flex items-center justify-center">
          <svg className="w-5 h-5 text-orange-600" fill="currentColor" viewBox="0 0 20 20">
            <path fillRule="evenodd" d="M3 5a2 2 0 012-2h10a2 2 0 012 2v8a2 2 0 01-2 2h-2.22l.123.489.804.804A1 1 0 0113 18H7a1 1 0 01-.707-1.707l.804-.804L7.22 15H5a2 2 0 01-2-2V5zm5.771 7H5V5h10v7H8.771z" clipRule="evenodd" />
          </svg>
        </div>
      ),
      title: 'Устройства',
      subtitle: '7',
      hasArrow: true
    }
  ];

  const renderMenuItem = (item: any) => (
    <button
      key={item.id}
      className="w-full flex items-center justify-between p-4 hover:bg-gray-50 transition-colors"
    >
      <div className="flex items-center space-x-3">
        {item.icon}
        <div className="text-left">
          <div className="text-base font-medium text-gray-900">{item.title}</div>
          {item.subtitle && (
            <div className="text-sm text-gray-500">{item.subtitle}</div>
          )}
        </div>
      </div>
      {item.hasArrow && (
        <svg className="w-5 h-5 text-gray-400" fill="currentColor" viewBox="0 0 20 20">
          <path fillRule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clipRule="evenodd" />
        </svg>
      )}
    </button>
  );

  return (
    <div className="min-h-screen bg-gray-50 pb-20">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 px-4 py-4 sticky top-0 z-10">
        <div className="flex items-center justify-between">
          <button className="p-2">
            <svg className="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
            </svg>
          </button>
          <h1 className="text-lg font-semibold text-gray-900">Профиль</h1>
          <button className="text-blue-600 font-medium">
            Изменить
          </button>
        </div>
      </header>

      {/* Profile Section */}
      <div className="bg-white">
        <div className="px-4 py-8 text-center">
          {/* Avatar */}
          <div className="relative inline-block mb-4">
            <img
              src={user.avatar}
              alt={user.name}
              className="w-24 h-24 rounded-full object-cover border-4 border-white shadow-lg"
            />
            <div className="absolute -bottom-1 -right-1 w-8 h-8 bg-orange-500 rounded-full flex items-center justify-center border-2 border-white">
              <span className="text-white text-xs font-bold">🔥</span>
            </div>
          </div>

          {/* User Info */}
          <h2 className="text-2xl font-bold text-gray-900 mb-1">{user.name}</h2>
          <p className="text-gray-600 mb-1">{user.phone}</p>
          <p className="text-blue-600">{user.username}</p>
        </div>

        {/* Profile Actions */}
        <div className="border-t border-gray-100">
          {menuItems.map(renderMenuItem)}
        </div>
      </div>

      {/* Support Section */}
      <div className="bg-white mt-4">
        {supportItems.map(renderMenuItem)}
      </div>

      {/* Add Account Section */}
      <div className="bg-white mt-4">
        {accountItems.map(renderMenuItem)}
      </div>

      {/* Other Options */}
      <div className="bg-white mt-4">
        {otherItems.map(renderMenuItem)}
      </div>

      <TabBar />
    </div>
  );
} 