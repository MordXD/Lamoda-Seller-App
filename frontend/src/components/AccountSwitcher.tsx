import { Fragment } from 'react';
import { Menu, Transition } from '@headlessui/react';
import { ChevronDownIcon, PlusIcon, UserIcon } from '@heroicons/react/24/outline';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function AccountSwitcher() {
  const { accounts, activeAccount, switchAccount, logout } = useAuth();
  const navigate = useNavigate();

  if (!activeAccount) {
    return null;
  }

  const handleAddAccount = () => {
    navigate('/login');
  };

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <Menu as="div" className="relative inline-block text-left">
      <div>
        <Menu.Button className="inline-flex w-full justify-center items-center gap-x-2 rounded-lg bg-white px-4 py-2 text-sm font-medium text-gray-700 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500">
          <UserIcon className="h-5 w-5 text-gray-400" />
          <span className="max-w-48 truncate">{activeAccount.shopName}</span>
          <ChevronDownIcon className="h-5 w-5 text-gray-400" />
        </Menu.Button>
      </div>

      <Transition
        as={Fragment}
        enter="transition ease-out duration-100"
        enterFrom="transform opacity-0 scale-95"
        enterTo="transform opacity-100 scale-100"
        leave="transition ease-in duration-75"
        leaveFrom="transform opacity-100 scale-100"
        leaveTo="transform opacity-0 scale-95"
      >
        <Menu.Items className="absolute right-0 z-10 mt-2 w-64 origin-top-right divide-y divide-gray-100 rounded-lg bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
          {/* Текущий аккаунт */}
          <div className="px-1 py-1">
            <div className="px-3 py-2 text-xs text-gray-500 uppercase tracking-wide">
              Текущий аккаунт
            </div>
            <div className="px-3 py-2 bg-blue-50 rounded-md">
              <div className="font-medium text-gray-900">{activeAccount.shopName}</div>
              <div className="text-sm text-gray-500">{activeAccount.user?.email}</div>
            </div>
          </div>

          {/* Другие аккаунты */}
          {accounts.length > 1 && (
            <div className="px-1 py-1">
              <div className="px-3 py-2 text-xs text-gray-500 uppercase tracking-wide">
                Переключить аккаунт
              </div>
              {accounts
                .filter(account => account.id !== activeAccount.id)
                .map((account) => (
                  <Menu.Item key={account.id}>
                    {({ active }) => (
                      <button
                        onClick={() => switchAccount(account.id)}
                        className={`${
                          active ? 'bg-gray-100' : ''
                        } group flex w-full items-center rounded-md px-3 py-2 text-sm text-gray-900 hover:bg-gray-100 transition-colors`}
                      >
                        <UserIcon className="mr-3 h-5 w-5 text-gray-400" />
                        <div className="text-left">
                          <div className="font-medium">{account.shopName}</div>
                          <div className="text-xs text-gray-500">{account.user?.email}</div>
                        </div>
                      </button>
                    )}
                  </Menu.Item>
                ))}
            </div>
          )}

          {/* Действия */}
          <div className="px-1 py-1">
            <Menu.Item>
              {({ active }) => (
                <button
                  onClick={handleAddAccount}
                  className={`${
                    active ? 'bg-gray-100' : ''
                  } group flex w-full items-center rounded-md px-3 py-2 text-sm text-gray-900 hover:bg-gray-100 transition-colors`}
                >
                  <PlusIcon className="mr-3 h-5 w-5 text-gray-400" />
                  Добавить аккаунт
                </button>
              )}
            </Menu.Item>
            <Menu.Item>
              {({ active }) => (
                <button
                  onClick={handleLogout}
                  className={`${
                    active ? 'bg-red-50 text-red-700' : 'text-red-600'
                  } group flex w-full items-center rounded-md px-3 py-2 text-sm hover:bg-red-50 hover:text-red-700 transition-colors`}
                >
                  Выйти
                </button>
              )}
            </Menu.Item>
          </div>
        </Menu.Items>
      </Transition>
    </Menu>
  );
} 