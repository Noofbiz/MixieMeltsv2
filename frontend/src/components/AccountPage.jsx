import { useAuth } from "../context/Context";
import AddProductForm from "./AddProductForm";
import AddSubscriptionBoxForm from "./AddSubscriptionBoxForm";

const mockOrders = [
  {
    id: "ORD-123",
    date: "2025-10-15",
    total: 12.49,
    status: "Delivered",
    items: [
      {
        name: "Lavender Dreams",
        quantity: 1,
        image:
          "https://placehold.co/400x400/e9d5ff/581c87?text=Lavender+Dreams",
      },
      {
        name: "Ocean Breeze",
        quantity: 1,
        image: "https://placehold.co/400x400/a5f3fc/155e75?text=Ocean+Breeze",
      },
    ],
  },
  {
    id: "ORD-456",
    date: "2025-09-22",
    total: 7.5,
    status: "Delivered",
    items: [
      {
        name: "Cozy Cashmere",
        quantity: 1,
        image: "https://placehold.co/400x400/e5e7eb/4b5563?text=Cozy+Cashmere",
      },
    ],
  },
  {
    id: "ORD-789",
    date: "2025-10-18",
    total: 31.48,
    status: "Shipped",
    items: [
      {
        name: "Spiced Apple Cider",
        quantity: 1,
        image: "https://placehold.co/400x400/fed7aa/9a3412?text=Spiced+Apple",
      },
      {
        name: "Monthly Subscription Box",
        quantity: 1,
        image:
          "https://placehold.co/400x400/dcfce7/166534?text=Subscription+Box",
      },
    ],
  },
];

const mockSubscriptions = [
  {
    id: "SUB-A1B",
    product: { name: "Monthly Subscription Box" },
    status: "Active",
    nextBilling: "2025-11-01",
  },
];

const AccountPage = () => {
  const { user } = useAuth();

  if (!user) {
    return (
      <div className="text-center py-20">
        Please log in to see your account details.
      </div>
    );
  }

  return (
    <div className="container mx-auto px-6 py-12">
      <h1 className="text-4xl font-bold font-serif text-brown-900 mb-8">
        Hello, {user.email}!
      </h1>
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Left Column: Settings */}
        <div className="lg:col-span-1 space-y-8">
          <div className="bg-white p-6 rounded-xl shadow-lg">
            <h2 className="text-2xl font-serif font-bold text-brown-900 mb-4">
              Account Settings
            </h2>
            <div className="space-y-3">
              <div>
                <label className="text-sm font-medium text-brown-800">
                  Email Address
                </label>
                <p className="text-brown-700">{user.email}</p>
              </div>
              <button className="w-full text-left px-4 py-2 text-sm text-terracotta-600 border border-terracotta-600 rounded-full hover:bg-terracotta-600 hover:text-cream-100 transition-colors">
                Change Password
              </button>
              <div className="pt-2">
                <label className="flex items-center space-x-2">
                  <input
                    type="checkbox"
                    className="rounded text-terracotta-500 focus:ring-terracotta-500"
                    defaultChecked
                  />
                  <span className="text-brown-800">
                    Receive promotional emails
                  </span>
                </label>
              </div>
            </div>
          </div>
          <div className="bg-white p-6 rounded-xl shadow-lg">
            <h2 className="text-2xl font-serif font-bold text-brown-900 mb-4">
              My Subscriptions
            </h2>
            {mockSubscriptions.map((sub) => (
              <div key={sub.id} className="border-t border-cream-200 pt-4 mt-4">
                <div className="flex justify-between items-center">
                  <p className="font-semibold text-brown-800">
                    {sub.product.name}
                  </p>
                  <span
                    className={`px-2 py-1 text-xs rounded-full ${sub.status === "Active" ? "bg-green-100 text-green-800" : "bg-gray-100 text-gray-800"}`}
                  >
                    {sub.status}
                  </span>
                </div>
                <p className="text-sm text-brown-700 mt-1">
                  Next billing date: {sub.nextBilling}
                </p>
                <button className="mt-2 text-sm text-red-500 hover:text-red-700">
                  Cancel Subscription
                </button>
              </div>
            ))}
          </div>
        </div>

        {/* Right Column: Order History */}
        <div className="lg:col-span-2 bg-white p-6 rounded-xl shadow-lg">
          <h2 className="text-2xl font-serif font-bold text-brown-900 mb-4">
            Order History
          </h2>
          <div className="space-y-6">
            {mockOrders.map((order) => (
              <div key={order.id} className="border-b border-cream-200 pb-4">
                <div className="flex justify-between items-center mb-2">
                  <div>
                    <p className="font-bold text-lg text-brown-900">
                      Order #{order.id}
                    </p>
                    <p className="text-sm text-brown-700">Date: {order.date}</p>
                  </div>
                  <div className="text-right">
                    <p className="font-bold text-lg text-brown-800">
                      ${order.total.toFixed(2)}
                    </p>
                    <p className="text-sm font-semibold text-brown-700">
                      {order.status}
                    </p>
                  </div>
                </div>
                {order.items.map((item) => (
                  <div
                    key={item.id}
                    className="flex items-center space-x-3 mt-2 pl-4"
                  >
                    <img
                      src={item.image}
                      alt={item.name}
                      className="w-12 h-12 object-cover rounded-md"
                    />
                    <div>
                      <p className="text-sm font-semibold text-brown-800">
                        {item.name}
                      </p>
                      <p className="text-xs text-brown-700">
                        Quantity: {item.quantity}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            ))}
          </div>
        </div>

        {user.is_admin && (
          <>
            <div className="lg:col-span-3">
              <AddProductForm />
            </div>
            <div className="lg:col-span-3">
              <AddSubscriptionBoxForm />
            </div>
          </>
        )}
      </div>
    </div>
  );
};

export default AccountPage;
