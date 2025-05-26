import React from 'react';
import { Globe, Shield, Zap, Settings } from 'lucide-react';

const features = [
  {
    icon: <Globe className="h-8 w-8 text-indigo-500" />,
    title: 'Decentralized Hosting',
    description: 'Your content is served from a decentralized network, making it censorship-resistant and always available.'
  },
  {
    icon: <Shield className="h-8 w-8 text-indigo-500" />,
    title: 'Enhanced Security',
    description: 'Web3 domains provide an extra layer of security through blockchain technology and cryptographic verification.'
  },
  {
    icon: <Zap className="h-8 w-8 text-indigo-500" />,
    title: 'Lightning Fast',
    description: 'Content delivery optimized for speed with global distribution and efficient caching mechanisms.'
  },
  {
    icon: <Settings className="h-8 w-8 text-indigo-500" />,
    title: 'Full Control',
    description: 'You have complete ownership of your domain with no intermediaries or central points of control.'
  }
];

const Features = () => {
  return (
    <section id="features" className="py-24 bg-white">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8">
        <div className="max-w-3xl mx-auto text-center mb-16">
          <h2 className="text-3xl sm:text-4xl font-bold text-gray-900 mb-4">
            Web3 Domain Features
          </h2>
          <p className="text-lg text-gray-600">
            Experience the benefits of decentralized web hosting with your new domain
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
          {features.map((feature, index) => (
            <div 
              key={index}
              className="relative p-8 rounded-xl bg-white border border-gray-100 shadow-sm hover:shadow-md transition-all duration-300 group"
            >
              <div className="absolute -top-4 left-8 p-3 rounded-lg bg-indigo-50 group-hover:bg-indigo-100 transition-colors duration-300">
                {feature.icon}
              </div>
              <h3 className="mt-6 text-xl font-semibold text-gray-900 mb-3">
                {feature.title}
              </h3>
              <p className="text-gray-600">
                {feature.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};

export default Features;