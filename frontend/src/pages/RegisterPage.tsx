import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

export default function RegisterPage() {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [phone, setPhone] = useState('');
  const [legalForm, setLegalForm] = useState('');
  const [companyName, setCompanyName] = useState('');
  const [cooperationModel, setCooperationModel] = useState('');
  const [productCategory, setProductCategory] = useState('');
  const [brandName, setBrandName] = useState('');
  const [trademarkDocument, setTrademarkDocument] = useState('');
  const [retailPresence, setRetailPresence] = useState('');
  const [websiteLink, setWebsiteLink] = useState('');
  const [skuCount, setSkuCount] = useState('');
  const [avgPrice, setAvgPrice] = useState('');
  const [brandOriginCountry, setBrandOriginCountry] = useState('');
  const [warehouseLocation, setWarehouseLocation] = useState('');
  const [priceListLink, setPriceListLink] = useState('');
  const [personalDataAgreement, setPersonalDataAgreement] = useState(false);
  const [confidentialityAgreement, setConfidentialityAgreement] = useState(false);

  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [generatedPassword, setGeneratedPassword] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (!name || !email) {
      setError('Пожалуйста, заполните обязательные поля: Имя и E-mail');
      return;
    }
    
    if (!personalDataAgreement || !confidentialityAgreement) {
      setError('Необходимо принять условия соглашений');
      return;
    }

    setIsLoading(true);

    try {
      // Отправляем только name и email, как и требовалось
      const response = await fetch('/auth/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name, email }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || 'Ошибка регистрации');
      }

      setGeneratedPassword(data.password);
    } catch (error: any) {
      console.error('Registration failed:', error);
      setError(error.message || 'Что-то пошло не так');
    } finally {
      setIsLoading(false);
    }
  };

  const handleCopyToClipboard = () => {
    navigator.clipboard.writeText(generatedPassword);
    alert('Пароль скопирован в буфер обмена!');
  };

  if (generatedPassword) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center px-4 font-montserrat-regular">
        <div className="max-w-md w-full bg-white rounded-2xl shadow-sm p-8 text-center">
          <h2 className="text-2xl font-bold mb-4">Регистрация успешна!</h2>
          <p className="mb-4">Ваш пароль для входа:</p>
          <div className="relative p-4 bg-gray-100 rounded-lg mb-4">
            <p className="text-lg font-mono break-all">{generatedPassword}</p>
            <button
              onClick={handleCopyToClipboard}
              className="absolute top-2 right-2 text-gray-500 hover:text-gray-800"
              title="Скопировать пароль"
            >
              <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
            </button>
          </div>
          <p className="text-sm text-gray-500 mb-6">Сохраните его в надежном месте. Вы можете использовать его для входа в свой аккаунт.</p>
          <button
            onClick={() => navigate('/login')}
            className="w-full bg-black text-white py-4 px-6 text-lg font-montserrat-semibold rounded-lg hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 transition-colors"
          >
            Перейти на страницу входа
          </button>
        </div>
      </div>
    );
  }

  const formInputClass = "w-full px-0 py-3 text-lg bg-transparent border-0 border-b-2 border-gray-200 focus:border-gray-400 focus:outline-none focus:ring-0 placeholder-gray-400 font-montserrat-regular";
  const formLabelClass = "block text-xs text-gray-500 mb-2 uppercase tracking-wide font-montserrat-regular";
  const formSelectClass = `${formInputClass} border-b-2`;

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4 font-montserrat-regular">
      <div className="max-w-4xl w-full bg-white rounded-2xl shadow-sm p-8">
        <div className="text-center">
          <img
            className="h-[40px] w-[220px] mx-auto mb-4"
            src="/icons/lamoda-icon.svg"
            alt="Регистрация"
          />
          <h2 className="text-3xl font-bold text-gray-900">
            Анкета для продавцов маркетплейса
          </h2>
        </div>

        <form onSubmit={handleSubmit} className="mt-8 space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-x-8 gap-y-6">
            {/* Ваше имя */}
            <div>
              <label htmlFor="name" className={formLabelClass}>
                Ваше имя <span className="text-red-500">*</span>
              </label>
              <input id="name" type="text" required value={name} onChange={(e) => setName(e.target.value)} className={formInputClass} placeholder="Иван" disabled={isLoading} />
            </div>

            {/* E-mail */}
            <div>
              <label htmlFor="email" className={formLabelClass}>
                E-mail <span className="text-red-500">*</span>
              </label>
              <input id="email" type="email" required value={email} onChange={(e) => setEmail(e.target.value)} className={formInputClass} placeholder="ivan@example.com" disabled={isLoading} />
            </div>

            {/* Телефон */}
            <div>
              <label htmlFor="phone" className={formLabelClass}>
                Телефон <span className="text-red-500">*</span>
              </label>
              <input id="phone" type="tel" required value={phone} onChange={(e) => setPhone(e.target.value)} className={formInputClass} placeholder="+7 (999) 999-99-99" disabled={isLoading} />
            </div>

            {/* Форма юридического лица */}
            <div>
              <label htmlFor="legalForm" className={formLabelClass}>
                Форма юридического лица <span className="text-red-500">*</span>
              </label>
              <select id="legalForm" required value={legalForm} onChange={(e) => setLegalForm(e.target.value)} className={formSelectClass} disabled={isLoading}>
                <option value="" disabled>Выберите форму</option>
                <option>ООО</option>
                <option>ОАО</option>
                <option>ЗАО</option>
                <option>ПАО</option>
                <option>АО</option>
                <option>ИП</option>
                <option>Не зарегистрировано</option>
              </select>
            </div>

            {/* Наименование компании */}
            <div className="md:col-span-2">
              <label htmlFor="companyName" className={formLabelClass}>
                Наименование компании <span className="text-red-500">*</span>
              </label>
              <input id="companyName" type="text" required value={companyName} onChange={(e) => setCompanyName(e.target.value)} className={formInputClass} placeholder="ООО Ромашка" disabled={isLoading} />
            </div>

            {/* Модель сотрудничества */}
            <div>
                <label htmlFor="cooperationModel" className={formLabelClass}>
                    Какая модель сотрудничества вам интересна <span className="text-red-500">*</span>
                </label>
                <select id="cooperationModel" required value={cooperationModel} onChange={(e) => setCooperationModel(e.target.value)} className={formSelectClass} disabled={isLoading}>
                    <option value="" disabled>Выберите модель</option>
                    <option>С хранением товара на складе Lamoda (FBO)</option>
                    <option>Отгрузки заказов со склада бренда, доставка до клиентов силами Lamoda (FBS)</option>
                    <option>Обе модели</option>
                </select>
            </div>

            {/* Категория товара */}
            <div>
                <label htmlFor="productCategory" className={formLabelClass}>
                    Категория товара <span className="text-red-500">*</span>
                </label>
                <select id="productCategory" required value={productCategory} onChange={(e) => setProductCategory(e.target.value)} className={formSelectClass} disabled={isLoading}>
                    <option value="" disabled>Выберите категорию</option>
                    <option>Одежда мужская</option>
                    <option>Одежда женская</option>
                    <option>Одежда детская</option>
                    <option>Одежда больших размеров</option>
                    <option>Одежда и аксессуары для беременных</option>
                    <option>Обувь мужская</option>
                    <option>Обувь женская</option>
                    <option>Обувь детская</option>
                    <option>Аксессуары и бижутерия</option>
                    <option>Детские аксессуары</option>
                    <option>Игрушки</option>
                    <option>Товары для дома</option>
                    <option>Косметика/Парфюмерия/Аксессуары для красоты</option>
                    <option>Багаж</option>
                    <option>Товары для спорта</option>
                    <option>Техника для красоты и здоровья</option>
                    <option>Ювелирные украшения</option>
                    <option>Книги</option>
                    <option>Premium</option>
                    <option>Другое</option>
                </select>
            </div>
            
            {/* Наименование бренда/брендов */}
            <div className="md:col-span-2">
              <label htmlFor="brandName" className={formLabelClass}>
                Наименование бренда/брендов <span className="text-red-500">*</span>
              </label>
              <input id="brandName" type="text" required value={brandName} onChange={(e) => setBrandName(e.target.value)} className={formInputClass} placeholder="Brand1, Brand2" disabled={isLoading} />
            </div>

            {/* Документ, подтверждающий право на Товарные знаки */}
            <div className="md:col-span-2">
                <label htmlFor="trademarkDocument" className={formLabelClass}>
                    Документ, подтверждающий право на Товарные знаки <span className="text-red-500">*</span>
                </label>
                <select id="trademarkDocument" required value={trademarkDocument} onChange={(e) => setTrademarkDocument(e.target.value)} className={formSelectClass} disabled={isLoading}>
                    <option value="" disabled>Выберите документ</option>
                    <option>Свидетельство на товарный знак (если Вы являетесь правообладателем)</option>
                    <option>Лицензионные/сублицензионные договоры (зарегистрированные)</option>
                    <option>Дистрибьюторские соглашения</option>
                    <option>Разрешительное письмо от правообладателя</option>
                    <option>Заявка на регистрацию товарного знака с отметкой уполномоченного органа о ее принятии к рассмотрению</option>
                </select>
            </div>
            
            {/* Есть ли собственная розница? */}
            <div className="md:col-span-2">
              <label htmlFor="retailPresence" className={formLabelClass}>
                Есть ли собственная розница? В каких розничных сетях вы представлены? <span className="text-red-500">*</span>
              </label>
              <textarea id="retailPresence" required value={retailPresence} onChange={(e) => setRetailPresence(e.target.value)} className={formInputClass} placeholder="Да, сеть 'Мода'..." disabled={isLoading} />
            </div>

            {/* Ссылка на сайт/социальные сети/каталог */}
            <div className="md:col-span-2">
              <label htmlFor="websiteLink" className={formLabelClass}>
                Ссылка на сайт/социальные сети/каталог <span className="text-red-500">*</span>
              </label>
              <input id="websiteLink" type="url" required value={websiteLink} onChange={(e) => setWebsiteLink(e.target.value)} className={formInputClass} placeholder="https://example.com" disabled={isLoading} />
            </div>

            {/* Количество артикулов */}
            <div>
              <label htmlFor="skuCount" className={formLabelClass}>
                Количество артикулов (SKU) <span className="text-red-500">*</span>
              </label>
              <input id="skuCount" type="text" required value={skuCount} onChange={(e) => setSkuCount(e.target.value)} className={formInputClass} placeholder="50" disabled={isLoading} />
            </div>

            {/* Средняя розничная цена */}
            <div>
              <label htmlFor="avgPrice" className={formLabelClass}>
                Средняя розничная цена, ₽ <span className="text-red-500">*</span>
              </label>
              <input id="avgPrice" type="text" required value={avgPrice} onChange={(e) => setAvgPrice(e.target.value)} className={formInputClass} placeholder="1500" disabled={isLoading} />
            </div>

            {/* Страна происхождения бренда */}
            <div>
              <label htmlFor="brandOriginCountry" className={formLabelClass}>
                Страна происхождения бренда <span className="text-red-500">*</span>
              </label>
              <input id="brandOriginCountry" type="text" required value={brandOriginCountry} onChange={(e) => setBrandOriginCountry(e.target.value)} className={formInputClass} placeholder="Россия" disabled={isLoading} />
            </div>

            {/* Местонахождение склада */}
            <div>
              <label htmlFor="warehouseLocation" className={formLabelClass}>
                Местонахождение склада <span className="text-red-500">*</span>
              </label>
              <input id="warehouseLocation" type="text" required value={warehouseLocation} onChange={(e) => setWarehouseLocation(e.target.value)} className={formInputClass} placeholder="Москва" disabled={isLoading} />
            </div>
            
            {/* Прайс-лист */}
            <div className="md:col-span-2">
              <label htmlFor="priceListLink" className={formLabelClass}>
                Ссылка на прайс-лист <span className="text-red-500">*</span>
              </label>
              <input id="priceListLink" type="url" required value={priceListLink} onChange={(e) => setPriceListLink(e.target.value)} className={formInputClass} placeholder="https://docs.google.com/spreadsheets/..." disabled={isLoading} />
              <p className="text-xs text-gray-500 mt-1">
                Заполните <a href="https://clck.ru/dXFmT" target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:underline">шаблон</a>, выложите на файлообменник и вставьте ссылку.
              </p>
            </div>
          </div>
          
          {/* Соглашения */}
          <div className="space-y-4">
            <div className="flex items-start">
              <input id="personal-data" name="personal-data" type="checkbox" required checked={personalDataAgreement} onChange={(e) => setPersonalDataAgreement(e.target.checked)} className="h-4 w-4 text-black border-gray-300 rounded focus:ring-black mt-1" />
              <div className="ml-3 text-sm">
                <label htmlFor="personal-data" className="font-medium text-gray-700">
                  Я даю согласие ООО "Купишуз" на обработку моих персональных данных в соответствии с <a href="https://www.lamoda.ru/about/personaldata/" target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:underline">Политикой обработки персональных данных</a> <span className="text-red-500">*</span>
                </label>
              </div>
            </div>
            <div className="flex items-start">
              <input id="confidentiality" name="confidentiality" type="checkbox" required checked={confidentialityAgreement} onChange={(e) => setConfidentialityAgreement(e.target.checked)} className="h-4 w-4 text-black border-gray-300 rounded focus:ring-black mt-1" />
              <div className="ml-3 text-sm">
                <label htmlFor="confidentiality" className="font-medium text-gray-700">
                  Я принимаю <a href="https://www.lamoda.ru/about/privacypolicy/" target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:underline">Соглашение о конфиденциальности</a> ООО "Купишуз" <span className="text-red-500">*</span>
                </label>
              </div>
            </div>
          </div>

          {error && (
            <div className="pt-4 text-red-600 text-sm text-center font-montserrat-regular">{error}</div>
          )}

          <div className="pt-4">
            <button
              type="submit"
              disabled={isLoading || !personalDataAgreement || !confidentialityAgreement}
              className="w-full bg-black text-white py-4 px-6 text-lg font-montserrat-semibold rounded-lg hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {isLoading ? 'Регистрация...' : 'Зарегистрироваться'}
            </button>
          </div>

          <div className="flex items-center justify-center pt-6">
            <button
              type="button"
              className="text-gray-600 hover:text-gray-800 text-base underline transition-colors font-montserrat-regular"
              onClick={() => navigate('/login')}
            >
              Уже есть аккаунт? Войти
            </button>
          </div>
        </form>
      </div>
    </div>
  );
} 